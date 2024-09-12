package universum

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestNewClient_Success(t *testing.T) {
	opts := &Options{
		ConnPoolsize: 10,
	}

	client, err := NewClient(opts)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.pool == nil {
		t.Fatal("Expected connection pool to be initialized, got nil")
	}

	if client.opts != opts {
		t.Fatalf("Expected client opts to match input, got %v", client.opts)
	}

	if client.id == "" {
		t.Fatal("Expected client ID to be set, but got an empty ID")
	}

	decodedId, err := base64.StdEncoding.DecodeString(client.id)
	if err != nil {
		t.Fatalf("Expected valid base64 ID, but decoding failed: %v", err)
	}

	_, err = strconv.Atoi(string(decodedId))
	if err != nil {
		t.Fatalf("Expected client ID to be a valid timestamp, but got: %v", err)
	}
}

func TestNewClient_InvalidPoolSize(t *testing.T) {
	opts := &Options{
		ConnPoolsize: 0,
	}

	client, err := NewClient(opts)

	if err == nil {
		t.Fatal("Expected error due to invalid pool size, got nil")
	}

	if client != nil {
		t.Fatalf("Expected client to be nil due to invalid pool size, but got %v", client)
	}
}

func TestClient_Commands(t *testing.T) {
	opts := mockOptions()

	client, err := NewClient(opts)
	if err != nil {
		t.Fatalf("Expected no error while creating client, got %v", err)
	}

	ctx := context.Background()

	var cError error

	_, cError = client.Info(ctx)
	if cError != nil {
		t.Fatalf("Expected no error from Get, got %v", err)
	}
}

// run this test only with a running universum server
func BenchmarkClient_Benchmark(b *testing.B) {
	opts := mockOptions()

	client, _ := NewClient(opts)
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("K-%v", time.Now().UnixNano())
		_, _ = client.Exists(ctx, key)
		_, _ = client.Set(ctx, key, 1034, 1000)
		_, _ = client.Get(ctx, key)
		_, _ = client.Increment(ctx, key, 12)
		_, _ = client.Decrement(ctx, key, 22)
		_, _ = client.MSet(ctx, map[string]interface{}{key + "AA": 200, key + "BB": "USD"})
		_, _ = client.MGet(ctx, []string{key + "AA", key + "BB", "CC"})
		//	_, _ = client.MDelete(ctx, []string{key + "AA", key + "BB"})
		_, _ = client.TTL(ctx, key)
		//	_, _ = client.Delete(ctx, key)
		_, _ = client.Ping(ctx)
		//_, _ = client.Info(ctx)
	}
}
