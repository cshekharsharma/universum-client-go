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

func TestClient_Commands(t *testing.T) {
	opts := mockOptions()

	client, err := NewClient(opts)
	if err != nil {
		t.Fatalf("Expected no error while creating client, got %v", err)
	}

	ctx := context.Background()
	var cError error

	info, cError := client.Info(ctx)
	if cError != nil {
		t.Fatalf("Expected no error from Get, got %v", err)
	}
	t.Logf("Exists: %v\n", info)

	key := "abcd"
	suffixAA := key + "AA"
	suffixBB := key + "BB"
	suffixCC := "CC"

	res1, err := client.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error from Exists, got %v\n", err)
	}
	t.Logf("Exists: %v\n", res1)

	res2, err := client.Set(ctx, key, 1034, 1000)
	if err != nil {
		t.Fatalf("Expected no error from Set, got %v\n", err)
	}
	t.Logf("Set: %v\n", res2)

	res3, err := client.Get(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error from Get, got %v\n", err)
	}
	t.Logf("Get: %v\n", res3)

	res4, err := client.Increment(ctx, key, 12)
	if err != nil {
		t.Fatalf("Expected no error from Increment, got %v\n", err)
	}
	t.Logf("Increment: %v\n", res4)

	res5, err := client.Decrement(ctx, key, 22)
	if err != nil {
		t.Fatalf("Expected no error from Decrement, got %v\n", err)
	}
	t.Logf("Decrement: %v\n", res5)

	res6, err := client.MSet(ctx, map[string]interface{}{suffixAA: []interface{}{10, 20}, suffixBB: "USD"})
	if err != nil {
		t.Fatalf("Expected no error from MSet, got %v\n", err)
	}
	t.Logf("MSet: %v\n", res6)

	res7, err := client.MGet(ctx, []string{suffixAA, suffixBB, suffixCC})
	if err != nil {
		t.Fatalf("Expected no error from MGet, got %v\n", err)
	}
	t.Logf("MGet: %#v\n", res7)

	res8, err := client.Append(ctx, suffixBB, "Only")
	if err != nil {
		t.Fatalf("Expected no error from Append, got %v\n", err)
	}
	t.Logf("Append: %v\n", res8)

	res9, err := client.MDelete(ctx, []string{suffixAA, suffixBB})
	if err != nil {
		t.Fatalf("Expected no error from MDelete, got %v\n", err)
	}
	t.Logf("MDelete: %v\n", res9)

	res10, err := client.TTL(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error from TTL, got %v\n", err)
	}
	t.Logf("TTL: %v\n", res10)

	res11, err := client.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error from Delete, got %v\n", err)
	}
	t.Logf("Delete: %v\n", res11)

	res12, err := client.Ping(ctx)
	if err != nil {
		t.Fatalf("Expected no error from Ping, got %v\n", err)
	}
	t.Logf("Ping: %v\n", res12)

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
		_, _ = client.MDelete(ctx, []string{key + "AA", key + "BB"})
		_, _ = client.TTL(ctx, key)
		_, _ = client.Delete(ctx, key)
		_, _ = client.Ping(ctx)
		_, _ = client.Info(ctx)
	}
}
