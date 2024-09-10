package universum

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var reusableTimers = sync.Pool{
	New: func() interface{} {
		t := time.NewTimer(time.Hour)
		t.Stop()
		return t
	},
}

type connPool struct {
	options   *Options
	connMutex sync.Mutex

	connections     []connInterface
	idleConnections []connInterface
	waitQueue       chan struct{}

	poolsize     int64
	numIdleConns int64
	isClosed     uint32
}

func (cp *connPool) createConn() (connInterface, error) {
	conn, err := newConnection(cp.options)
	if cp.poolsize < cp.options.ConnPoolsize {
		conn.setPooled(true)
		cp.poolsize++
	}
	return conn, err
}

func (cp *connPool) GetConn(ctx context.Context) (connInterface, error) {
	if cp.closed() {
		return nil, ErrConnectionPoolClosed
	}

	if err := cp.waitForTurn(ctx); err != nil {
		return nil, err
	}

	for {
		cp.connMutex.Lock()
		conn, err := cp.acquireIdleConnection()
		cp.connMutex.Unlock()

		if err != nil {
			cp.freeTurn()
			return nil, err
		}

		if conn == nil {
			break
		}

		if !cp.isActiveConnection(conn) {
			cp.CloseConn(conn)
			continue
		}

		return conn, nil
	}

	newConn, err := cp.createConn()
	if err != nil {
		cp.freeTurn()
		return nil, err
	}

	cp.connMutex.Lock()
	cp.connections = append(cp.connections, newConn)
	cp.connMutex.Unlock()

	return newConn, nil
}

func (cp *connPool) waitForTurn(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	select {
	case cp.waitQueue <- struct{}{}:
		return nil
	default:
	}

	timer := reusableTimers.Get().(*time.Timer)
	timer.Reset(cp.options.ConnWaitTimeout)

	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		reusableTimers.Put(timer)
		return ctx.Err()

	case cp.waitQueue <- struct{}{}:
		if !timer.Stop() {
			<-timer.C
		}
		reusableTimers.Put(timer)
		return nil

	case <-timer.C:
		reusableTimers.Put(timer)
		return fmt.Errorf("request timed out while waiting for turn in pool: %w", ErrConnectionWaitTimeout)
	}
}

func (cp *connPool) freeTurn() {
	<-cp.waitQueue
}

func (cp *connPool) isActiveConnection(conn connInterface) bool {
	now := time.Now()

	if cp.options.ConnMaxLifetime > 0 && now.Sub(conn.getCreatedAt()) >= cp.options.ConnMaxLifetime {
		return false
	}

	_ = conn.getNetConn().SetDeadline(time.Now().Add(1 * time.Second))

	sysConn, ok := conn.getNetConn().(syscall.Conn)
	if !ok {
		return true
	}

	rawConn, err := sysConn.SyscallConn()
	if err != nil {
		return false
	}

	var sysErr error

	if err := rawConn.Read(func(fd uintptr) bool {
		var buf [1]byte
		bytesRead, err := syscall.Read(int(fd), buf[:])

		switch {
		case bytesRead == 0 && err == nil:
			sysErr = io.EOF
		case bytesRead > 0:
			sysErr = errUnexpectedRead
		case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
			sysErr = nil
		default:
			sysErr = err
		}
		return sysErr == nil
	}); err != nil {
		return false
	}

	if sysErr != nil {
		return false
	}

	conn.setUsedAt(now)
	return true
}

func (cp *connPool) acquireIdleConnection() (connInterface, error) {
	if cp.closed() {
		return nil, ErrConnectionPoolClosed
	}

	idleQueueSize := len(cp.idleConnections)
	if idleQueueSize == 0 {
		return nil, nil
	}

	conn := cp.idleConnections[0]
	copy(cp.idleConnections, cp.idleConnections[1:])
	cp.idleConnections = cp.idleConnections[:idleQueueSize-1]

	cp.numIdleConns--
	return conn, nil
}

func (cp *connPool) ReleaseConn(ctx context.Context, conn connInterface) {
	if conn.getReader().Buffered() > 0 {
		cp.Remove(ctx, conn)
		return
	}

	if !conn.getPooled() {
		cp.Remove(ctx, conn)
		return
	}

	var shouldCloseConn bool

	cp.connMutex.Lock()

	if cp.numIdleConns < cp.options.ConnPoolsize {
		cp.idleConnections = append(cp.idleConnections, conn)
		cp.numIdleConns++
	} else {
		cp.removeConnFromPool(conn)
		shouldCloseConn = true
	}

	cp.connMutex.Unlock()

	cp.freeTurn()

	if shouldCloseConn {
		cp.closeConn(conn)
	}
}

func (cp *connPool) Remove(_ context.Context, conn connInterface) {
	cp.removeConnFromPoolWithLock(conn)
	cp.freeTurn()
	cp.closeConn(conn)
}

func (cp *connPool) CloseConn(conn connInterface) error {
	cp.removeConnFromPoolWithLock(conn)
	return cp.closeConn(conn)
}

func (cp *connPool) removeConnFromPoolWithLock(conn connInterface) {
	cp.connMutex.Lock()
	defer cp.connMutex.Unlock()
	cp.removeConnFromPool(conn)
}

func (cp *connPool) removeConnFromPool(conn connInterface) {
	for index, currConn := range cp.connections {
		if currConn == conn {
			cp.connections = append(cp.connections[:index], cp.connections[index+1:]...)
			if conn.getPooled() {
				cp.poolsize--
			}
			break
		}
	}
}

func (cp *connPool) closeConn(conn connInterface) error {
	return conn.close()
}

// Len returns total number of connections.
func (cp *connPool) Len() int {
	cp.connMutex.Lock()
	length := len(cp.connections)
	cp.connMutex.Unlock()
	return length
}

// IdleLen returns number of idle connections.
func (cp *connPool) IdleLen() int {
	cp.connMutex.Lock()
	length := cp.numIdleConns
	cp.connMutex.Unlock()
	return int(length)
}

func (cp *connPool) closed() bool {
	return atomic.LoadUint32(&cp.isClosed) == 1
}

func (cp *connPool) Close() error {
	if !atomic.CompareAndSwapUint32(&cp.isClosed, 0, 1) {
		return ErrConnectionPoolClosed
	}

	var firstErr error
	cp.connMutex.Lock()
	for _, conn := range cp.connections {
		if conn != nil {
			if err := cp.closeConn(conn); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}

	cp.connections = nil
	cp.poolsize = 0
	cp.idleConnections = nil
	cp.numIdleConns = 0
	cp.connMutex.Unlock()

	return firstErr
}

//////////////////////////////////////////////////////////////////////////////

func newConnPool(opts *Options) (*connPool, error) {
	if opts.ConnPoolsize <= 0 {
		return nil, errors.New("connection pool size must be greater than 0")
	}

	pool := &connPool{
		options:         opts,
		connMutex:       sync.Mutex{},
		connections:     make([]connInterface, 0, opts.ConnPoolsize),
		idleConnections: make([]connInterface, 0, opts.ConnPoolsize),
		waitQueue:       make(chan struct{}, opts.ConnPoolsize*2),
		poolsize:        0,
		numIdleConns:    0,
		isClosed:        0,
	}

	return pool, nil
}
