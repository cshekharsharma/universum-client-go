package universum

import (
	"errors"
	"io"
	"net"
	"syscall"
	"time"
)

var errUnexpectedRead = errors.New("unexpected read from socket")

type connPool struct {
}

func (cp *connPool) getConn() {

}

func newConnPool() *connPool {
	return &connPool{}
}

func checkConnection(conn net.Conn) error {
	_ = conn.SetDeadline(time.Time{})

	sysConn, ok := conn.(syscall.Conn)
	if !ok {
		return nil
	}
	rawConn, err := sysConn.SyscallConn()
	if err != nil {
		return err
	}

	var sysErr error

	if err := rawConn.Read(func(fd uintptr) bool {
		var buf [1]byte
		n, err := syscall.Read(int(fd), buf[:])
		switch {
		case n == 0 && err == nil:
			sysErr = io.EOF
		case n > 0:
			sysErr = errUnexpectedRead
		case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
			sysErr = nil
		default:
			sysErr = err
		}
		return true
	}); err != nil {
		return err
	}

	return sysErr
}
