package universum

import (
	"context"
	"errors"
	"testing"
	"time"
)

// MockOptions creates mock options for testing purposes
func mockOptions() *Options {
	return &Options{
		ConnPoolsize:    10,
		ConnWaitTimeout: 10 * time.Second,
		ConnMaxLifetime: 1 * time.Hour,
		HostAddr:        "localhost:11191",
		DialTimeout:     1 * time.Second,
		MaxRetries:      5,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
	}
}

// TestNewConnPool creates a new connection pool and verifies its initialization
func TestNewConnPool(t *testing.T) {
	opts := mockOptions()

	pool, err := newConnPool(opts)
	if err != nil {
		t.Fatalf("Expected to create conn pool, got error: %v", err)
	}

	if pool == nil {
		t.Fatalf("Expected non-nil pool")
	}

	if len(pool.connections) != int(pool.poolsize) {
		t.Errorf("Expected pool size %d, got %d", pool.poolsize, len(pool.connections))
	}
}

// TestGetConn acquires a connection from the pool
func TestGetConn(t *testing.T) {
	opts := mockOptions()
	connPool, _ := newConnPool(opts)

	ctx := context.Background()
	conn, err := connPool.GetConn(ctx)
	if err != nil {
		t.Fatalf("Expected to acquire connection, got error: %v", err)
	}

	if conn == nil {
		t.Fatalf("Expected non-nil connection, got nil")
	}
}

// TestReleaseConn verifies releasing a connection back to the pool
func TestReleaseConn(t *testing.T) {
	opts := mockOptions()
	pool, _ := newConnPool(opts)

	ctx := context.Background()
	conn, err := pool.GetConn(ctx)
	if err != nil {
		t.Fatalf("Failed to acquire connection: %v", err)
	}

	// Release connection back to the pool
	pool.ReleaseConn(ctx, conn)

	if pool.IdleLen() != 1 {
		t.Errorf("Expected idle connections to be 1, got %d", pool.IdleLen())
	}
}

// TestCloseConnPool tests the closing of the pool
func TestCloseConnPool(t *testing.T) {
	opts := mockOptions()
	pool, _ := newConnPool(opts)

	err := pool.Close()
	if err != nil {
		t.Fatalf("Expected to close connection pool without error, got: %v", err)
	}

	if !pool.closed() {
		t.Fatal("Expected pool to be closed")
	}
}

// TestWaitForTurnTimeout verifies that waitForTurn times out correctly
func TestWaitForTurnTimeout(t *testing.T) {
	opts := mockOptions()
	opts.ConnWaitTimeout = 500 * time.Millisecond
	pool, _ := newConnPool(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	time.Sleep(300 * time.Millisecond)
	err := pool.waitForTurn(ctx)

	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Expected context deadline exceeded error, got %v", err)
	}
}

// TestIsActiveConnection verifies if the connection is active
func TestIsActiveConnection(t *testing.T) {
	opts := mockOptions()
	pool, _ := newConnPool(opts)

	conn, _ := pool.GetConn(context.Background())

	conn.setCreatedAt(time.Now().Add(-15 * time.Minute))
	if !pool.isActiveConnection(conn) {
		t.Fatal("Expected connection to be inactive, but it is active")
	}
}
