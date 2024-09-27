package universum

import (
	"bufio"
	"context"
	"net"
	"testing"
	"time"
)

// MockNetConn implements net.Conn for testing purposes
type MockNetConn struct{}

func (m *MockNetConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (m *MockNetConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (m *MockNetConn) Close() error                       { return nil }
func (m *MockNetConn) LocalAddr() net.Addr                { return nil }
func (m *MockNetConn) RemoteAddr() net.Addr               { return nil }
func (m *MockNetConn) SetDeadline(t time.Time) error      { return nil }
func (m *MockNetConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *MockNetConn) SetWriteDeadline(t time.Time) error { return nil }

// TestNewConnection tests the creation of a new connection
func TestNewConnection(t *testing.T) {
	opts := &Options{
		HostAddr:    "localhost:11191",
		DialTimeout: 1 * time.Second,
		MaxRetries:  1,
	}

	opts.Init()
	_, err := newConnection(opts)
	if err != nil {
		t.Fatalf("Expected connection dial failed error, got %v", err)
	}
}

// TestWrite tests the write functionality
func TestConnWrite(t *testing.T) {
	netConn := &MockNetConn{}
	conn := &Conn{
		netconn: netConn,
		writer:  bufio.NewWriter(netConn),
	}

	n, err := conn.write([]byte("hello"))
	if err != nil {
		t.Fatalf("Unexpected error during write: %v", err)
	}

	if n != 5 {
		t.Errorf("Expected to write 5 bytes, wrote %d", n)
	}
}

// TestGetSetCreatedAt tests getting and setting the created time of the connection
func TestGetSetCreatedAt(t *testing.T) {
	conn := &Conn{}
	now := time.Now()

	conn.setCreatedAt(now)
	if conn.getCreatedAt() != now {
		t.Errorf("Expected createdAt to be %v, got %v", now, conn.getCreatedAt())
	}
}

// TestGetSetUsedAt tests getting and setting the used time of the connection
func TestGetSetUsedAt(t *testing.T) {
	conn := &Conn{}
	now := time.Now()

	conn.setUsedAt(now)
	if conn.getUsedAt() != now {
		if now.Unix() != conn.getUsedAt().Unix() {
			t.Errorf("Expected usedAt to be %v, got %v", now.Unix(), conn.getUsedAt().Unix())
		}
	}
}

// TestGetSetInUse tests setting and getting the in-use status of the connection
func TestGetSetInUse(t *testing.T) {
	conn := &Conn{}
	conn.setInUse(true)
	if !conn.getInUse() {
		t.Errorf("Expected connection to be in use, but it is not")
	}

	conn.setInUse(false)
	if conn.getInUse() {
		t.Errorf("Expected connection to not be in use, but it is")
	}
}

// TestClose tests closing the connection
func TestConnClose(t *testing.T) {
	netConn := &MockNetConn{}
	conn := &Conn{netconn: netConn}

	err := conn.close()
	if err != nil {
		t.Errorf("Unexpected error during close: %v", err)
	}

	if conn.getInUse() {
		t.Errorf("Expected connection to be not in use after close")
	}
}

// TestDeadline tests the deadline calculation
func TestDeadline(t *testing.T) {
	conn := &Conn{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	timeout := 1 * time.Second
	deadline := conn.deadline(ctx, timeout)

	expectedDeadline := time.Now().Add(timeout)
	if deadline.Truncate(time.Millisecond).Before(expectedDeadline.Truncate(time.Millisecond)) {
		t.Errorf("Expected deadline to be %v, got %v", expectedDeadline, deadline)
	}
}
