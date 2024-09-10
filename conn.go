package universum

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

const tcpDialer string = "tcp"

var noDeadline time.Time = time.Time{}

type connInterface interface {
	write(content []byte) (int, error)
	close() error

	getCreatedAt() time.Time
	getUsedAt() time.Time
	getInUse() bool
	getRemoteAddr() net.Addr
	getNetConn() net.Conn
	getPooled() bool
	getReader() *bufio.Reader
	getWriter() *bufio.Writer

	setUsedAt(t time.Time)
	setInUse(state bool)
	setPooled(pooled bool)
	setCreatedAt(time.Time)
}

// Conn represents a connection structure
type Conn struct {
	netconn   net.Conn
	writer    *bufio.Writer
	reader    *bufio.Reader
	pooled    bool
	createdAt time.Time
	usedAt    int64
	inUse     int32
}

// Write writes content to the connection's buffer
func (c *Conn) write(content []byte) (int, error) {
	return c.writer.Write(content)
}

// CreatedAt returns the time the connection was created
func (c *Conn) getCreatedAt() time.Time {
	return c.createdAt
}

// SetCreatedAt sets the time the connection was created
func (c *Conn) setCreatedAt(t time.Time) {
	c.createdAt = t
}

// GetUsedAt returns the last time the connection was used
func (c *Conn) getUsedAt() time.Time {
	unix := atomic.LoadInt64(&c.usedAt)
	return time.Unix(unix, 0)
}

// SetUsedAt sets the time when the connection was last used
func (c *Conn) setUsedAt(t time.Time) {
	atomic.StoreInt64(&c.usedAt, t.Unix())
}

// InUse returns whether the connection is currently in use
func (c *Conn) getInUse() bool {
	return atomic.LoadInt32(&c.inUse) == 1
}

// SetInUse safely sets the in-use status of the connection
func (c *Conn) setInUse(state bool) {
	if state {
		atomic.StoreInt32(&c.inUse, 1)
	} else {
		atomic.StoreInt32(&c.inUse, 0)
	}
}

// RemoteAddr returns the remote address of the connection
func (c *Conn) getRemoteAddr() net.Addr {
	return c.netconn.RemoteAddr()
}

// GetPooled returns whether the connection is pooled
func (c *Conn) getPooled() bool {
	return c.pooled
}

// SetPooled sets the connection as pooled or not
func (c *Conn) setPooled(pooled bool) {
	c.pooled = pooled
}

// GetNetConn returns the underlying net.Conn
func (c *Conn) getNetConn() net.Conn {
	return c.netconn
}

// GetReader returns the reader
func (c *Conn) getReader() *bufio.Reader {
	return c.reader
}

// GetWriter returns the writer
func (c *Conn) getWriter() *bufio.Writer {
	return c.writer
}

// Close closes the connection and marks it as no longer in use
func (c *Conn) close() error {
	c.setInUse(false)
	return c.netconn.Close()
}

// deadline calculates the appropriate deadline based on the context and timeout
func (c *Conn) deadline(ctx context.Context, timeout time.Duration) time.Time {
	tm := time.Now()
	c.setUsedAt(tm)

	if timeout > 0 {
		tm = tm.Add(timeout)
	}

	if ctx != nil {
		deadline, ok := ctx.Deadline()
		if ok {
			if timeout == 0 {
				return deadline
			}
			if deadline.Before(tm) {
				return deadline
			}
			return tm
		}
	}

	if timeout > 0 {
		return tm
	}

	return noDeadline
}

// newConnection creates a new connection to the specified address
func newConnection(opts *Options) (connInterface, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opts.DialTimeout*time.Second)
	defer cancel()

	var dialer net.Dialer = net.Dialer{}
	var retryCount int64 = 0
	var dialedConn net.Conn
	var connErr error

	for retryCount < opts.MaxRetries {
		retryCount++
		dialedConn, connErr = dialer.DialContext(ctx, tcpDialer, opts.HostAddr)

		if connErr != nil {
			dialedConn.Close()

			if ctx.Err() == context.DeadlineExceeded {
				connErr = fmt.Errorf("dial to host %s failed due to timeout after %s: %w",
					opts.HostAddr, opts.DialTimeout, ErrConnectionDialTimeout)
			} else {
				connErr = fmt.Errorf("failed to dial host %s [%v]: %w",
					opts.HostAddr, connErr, ErrConnectionDialFailed)
			}

			continue
		}
	}

	if connErr != nil {
		return nil, connErr
	}

	conn := &Conn{
		netconn:   dialedConn,
		reader:    bufio.NewReader(dialedConn),
		writer:    bufio.NewWriter(dialedConn),
		createdAt: time.Now(),
		pooled:    true,
	}

	var err error
	if opts.ReadTimeout > 0 {
		err = dialedConn.SetReadDeadline(time.Now().Add(opts.ReadTimeout))
		if err != nil {
			dialedConn.Close()
			return nil, fmt.Errorf("failed to set read deadline: %w", ErrConnectionConfigFailed)
		}
	}

	if opts.WriteTimeout > 0 {
		err = dialedConn.SetWriteDeadline(time.Now().Add(opts.WriteTimeout))
		if err != nil {
			dialedConn.Close()
			return nil, fmt.Errorf("failed to set write deadline: %w", ErrConnectionConfigFailed)
		}
	}

	conn.setUsedAt(time.Now())
	conn.setInUse(false)

	return conn, nil
}
