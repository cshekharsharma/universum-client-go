package universum

import (
	"bufio"
	"context"
	"net"
	"sync/atomic"
	"time"
)

const tcpDialer string = "tcp"

var noDeadline time.Time = time.Time{}

type Conn struct {
	netconn net.Conn

	writer  *bufio.Writer
	reader  *bufio.Reader
	wreader *bufio.Writer

	createdAt time.Time
	usedAt    int64
	inUse     bool
}

func (c *Conn) Write(content []byte) (int, error) {
	return c.writer.Write(content)
}

func (c *Conn) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Conn) GetUsedAt() time.Time {
	unix := atomic.LoadInt64(&c.usedAt)
	return time.Unix(unix, 0)
}

func (c *Conn) SetUsedAt(t time.Time) {
	atomic.StoreInt64(&c.usedAt, t.Unix())
}

func (c *Conn) InUsed() bool {
	return c.inUse
}

func (c *Conn) RemoteAdd() net.Addr {
	return c.netconn.RemoteAddr()
}

func (c *Conn) WithReader(
	ctx context.Context, timeout time.Duration, fn func(rd *bufio.Reader) error,
) error {
	if timeout >= 0 {
		if err := c.netconn.SetReadDeadline(c.deadline(ctx, timeout)); err != nil {
			return err
		}
	}
	return fn(c.reader)
}

func (c *Conn) WithWriter(
	ctx context.Context, timeout time.Duration, fn func(wr *bufio.Writer) error,
) error {
	if timeout >= 0 {
		if err := c.netconn.SetWriteDeadline(c.deadline(ctx, timeout)); err != nil {
			return err
		}
	}

	if c.writer.Buffered() > 0 {
		c.writer.Reset(c.netconn)
	}

	if err := fn(c.wreader); err != nil {
		return err
	}

	return c.writer.Flush()
}

func (c *Conn) Close() error {
	return c.netconn.Close()
}

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

func newConnection(addr string) (*Conn, error) {
	netconn, err := net.Dial(tcpDialer, addr)

	if err != nil {
		return nil, err
	}

	conn := &Conn{
		netconn:   netconn,
		reader:    bufio.NewReader(netconn),
		writer:    bufio.NewWriter(netconn),
		createdAt: time.Now(),
		inUse:     false,
	}

	conn.SetUsedAt(time.Now())
	return conn, nil
}
