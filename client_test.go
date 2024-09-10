package universum

import (
	"context"
	"encoding/base64"
	"strconv"
	"testing"
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

func TestClient_Get_Success(t *testing.T) {
	opts := mockOptions()

	client, err := NewClient(opts)
	if err != nil {
		t.Fatalf("Expected no error while creating client, got %v", err)
	}

	ctx := context.Background()
	key := "some_key"

	result, err := client.Get(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error from Get, got %v", err)
	}

	// Validate the result
	if result == nil {
		t.Fatal("Expected a valid CommandResult, but got nil")
	}

	if result.code != 200 || result.message != "OK" {
		t.Fatalf("Expected CommandResult with code 200 and message OK, got code %v and message %v", result.code, result.message)
	}
}

func TestClient_Get_CommandError(t *testing.T) {
	opts := mockOptions()

	client, err := NewClient(opts)
	if err != nil {
		t.Fatalf("Expected no error while creating client, got %v", err)
	}

	ctx := context.Background()
	key := "some_key"

	_, err = client.Get(ctx, key)
	if err == nil {
		t.Fatal("Expected error from Get due to command failure, but got nil")
	}

	if err.Error() != "command error" {
		t.Fatalf("Unexpected error message, got: %v", err)
	}
}

func TestClient_Get_InvalidCommandResult(t *testing.T) {
	opts := mockOptions()

	client, err := NewClient(opts)
	if err != nil {
		t.Fatalf("Expected no error while creating client, got %v", err)
	}

	ctx := context.Background()
	key := "some_key"

	_, err = client.Get(ctx, key)
	if err == nil {
		t.Fatal("Expected error from Get due to invalid command result, but got nil")
	}

	if err.Error() != "invalid result from server found: malformed response received" {
		t.Fatalf("Unexpected error message, got: %v", err)
	}
}
