package universum

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

const (
	tcpDialer string = "tcp"
)

var (
	noDeadline time.Time = time.Time{}
)

// Conn represents a connection structure
type Conn struct {
	netconn net.Conn

	writer *bufio.Writer
	reader *bufio.Reader

	Inited bool
	pooled bool

	createdAt time.Time
	usedAt    int64
	inUse     int32 // Use atomic int32 for thread-safe access
}

// Write writes content to the connection's buffer
func (c *Conn) Write(content []byte) (int, error) {
	return c.writer.Write(content)
}

// CreatedAt returns the time the connection was created
func (c *Conn) CreatedAt() time.Time {
	return c.createdAt
}

// GetUsedAt returns the last time the connection was used
func (c *Conn) GetUsedAt() time.Time {
	unix := atomic.LoadInt64(&c.usedAt)
	return time.Unix(unix, 0)
}

// SetUsedAt sets the time when the connection was last used
func (c *Conn) SetUsedAt(t time.Time) {
	atomic.StoreInt64(&c.usedAt, t.Unix())
}

// InUse returns whether the connection is currently in use
func (c *Conn) InUse() bool {
	return atomic.LoadInt32(&c.inUse) == 1
}

// SetInUse safely sets the in-use status of the connection
func (c *Conn) SetInUse(state bool) {
	if state {
		atomic.StoreInt32(&c.inUse, 1)
	} else {
		atomic.StoreInt32(&c.inUse, 0)
	}
}

// RemoteAdd returns the remote address of the connection
func (c *Conn) RemoteAddr() net.Addr {
	return c.netconn.RemoteAddr()
}

// WithReader safely executes a function with a reader, applying deadlines if necessary
func (c *Conn) WithReader(ctx context.Context, timeout time.Duration, fn func(rd *bufio.Reader) error) error {
	if deadline, ok := ctx.Deadline(); ok {
		// Use the context's deadline if set
		if err := c.netconn.SetReadDeadline(deadline); err != nil {
			return err
		}
	} else if timeout > 0 {
		// Otherwise, use the provided timeout
		if err := c.netconn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
			return err
		}
	}

	// Execute the function
	return fn(c.reader)
}

// WithWriter safely executes a function with a writer, applying deadlines if necessary
func (c *Conn) WithWriter(ctx context.Context, timeout time.Duration, fn func(wr *bufio.Writer) error) error {
	if deadline, ok := ctx.Deadline(); ok {
		if err := c.netconn.SetWriteDeadline(deadline); err != nil {
			return err
		}
	} else if timeout > 0 {
		if err := c.netconn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
			return err
		}
	}

	if c.writer.Buffered() > 0 {
		c.writer.Reset(c.netconn)
	}

	if err := fn(c.writer); err != nil {
		return err
	}

	return c.writer.Flush()
}

// Close closes the connection and marks it as no longer in use
func (c *Conn) Close() error {
	c.SetInUse(false)
	return c.netconn.Close()
}

// deadline calculates the appropriate deadline based on the context and timeout
func (c *Conn) deadline(ctx context.Context, timeout time.Duration) time.Time {
	tm := time.Now()
	c.SetUsedAt(tm)

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
func newConnection(opts *Options) (*Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opts.DialTimeout*time.Second)
	defer cancel()

	var dialer net.Dialer = net.Dialer{}
	var retryCount int64 = 0
	var netconn net.Conn
	var connErr error

	for retryCount <= opts.MaxRetries {
		netconn, err := dialer.DialContext(ctx, tcpDialer, opts.HostAddr)

		if err != nil {
			netconn.Close()

			if ctx.Err() == context.DeadlineExceeded {
				connErr = fmt.Errorf("dial to host %s failed due to timeout after %s: %w",
					opts.HostAddr, opts.DialTimeout, ErrConnectionDialTimeout)
			} else {
				connErr = fmt.Errorf("failed to dial host %s [%v]: %w",
					opts.HostAddr, err, ErrConnectionDialFailed)
			}

			retryCount++
			continue
		}
	}

	if connErr != nil {
		return nil, connErr
	}

	conn := &Conn{
		netconn:   netconn,
		reader:    bufio.NewReader(netconn),
		writer:    bufio.NewWriter(netconn),
		createdAt: time.Now(),
	}

	// Initialize the used time to the current time
	conn.SetUsedAt(time.Now())
	conn.SetInUse(false)

	return conn, nil
}
