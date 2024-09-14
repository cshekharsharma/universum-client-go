// Package universum provides a Go client for interacting with Universum DB.
// This client manages connections through a connection pool and allows you to perform
// operations like Get, Set, Delete, Increment, and others on key-value pairs in the database.
package universum

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var ncmu sync.Mutex

// Client is the main struct that represents a connection to the Universum NoSQL database.
// It includes a connection pool and holds client-specific options.
//
// Fields:
// - id: A unique identifier for the client, typically encoded as a base64 string.
// - pool: A pool of connections to manage database interactions.
// - opts: Configuration options provided to the client.
type Client struct {
	id   string
	pool *connPool
	opts *Options
}

// Get retrieves the value of a specified key from the Universum database.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
// - key: The key to retrieve the value for.
//
// Returns:
// - *GetResult: The result of the GET operation.
// - error: Returns an error if the command fails.
func (c *Client) Get(ctx context.Context, key string) (*GetResult, error) {
	result, err := sendCommand(ctx, c, commandGet, key)
	if err != nil {
		return nil, err
	}

	if decoded, ok := result.value.(map[string]interface{}); ok {
		return &GetResult{
			Value: decoded["Value"],
			Code:  result.code,
		}, nil
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// Set sets the value of a specified key in the Universum database with an optional TTL (time-to-live).
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
// - key: The key to set the value for.
// - value: The value to store under the given key.
// - ttl: The time-to-live for the key in seconds. If ttl is zero, the key will not expire.
//
// Returns:
// - *SetResult: The result of the SET operation.
// - error: Returns an error if the command fails.
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl int64) (*SetResult, error) {
	if c.opts.IsReadonly {
		return nil, fmt.Errorf("cannot execute write op in read-only client: %w", ErrClientReadonly)
	}

	if !isWriteableDatatype(value) {
		return nil, fmt.Errorf("provided datatype is not supported for write operations, "+
			"only int|float|bool|string|[]interface{} types are supported: %w", ErrInvalidDatatype)
	}

	result, err := sendCommand(ctx, c, commandSet, key, value, ttl)
	if err != nil {
		return nil, err
	}

	if didSet, ok := result.value.(bool); ok {
		return &SetResult{
			Success: didSet,
			Code:    result.code,
		}, nil
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// Exists checks if a specified key exists in the Universum database.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
// - key: The key to check for existence.
//
// Returns:
// - *ExistsResult: The result of the EXISTS operation.
// - error: Returns an error if the command fails.
func (c *Client) Exists(ctx context.Context, key string) (*ExistsResult, error) {
	result, err := sendCommand(ctx, c, commandExists, key)
	if err != nil {
		return nil, err
	}

	if found, ok := result.value.(bool); ok {
		return &ExistsResult{
			Found: found,
			Code:  result.code,
		}, nil
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// Delete removes the specified key from the Universum database.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
// - key: The key to delete.
//
// Returns:
// - *DeleteResult: The result of the DELETE operation.
// - error: Returns an error if the command fails.
func (c *Client) Delete(ctx context.Context, key string) (*DeleteResult, error) {
	if c.opts.IsReadonly {
		return nil, fmt.Errorf("cannot execute write op in read-only client: %w", ErrClientReadonly)
	}

	result, err := sendCommand(ctx, c, commandDelete, key)
	if err != nil {
		return nil, err
	}

	if deleted, ok := result.value.(bool); ok {
		return &DeleteResult{
			Deleted: deleted,
			Code:    result.code,
		}, nil
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// Increment increases the value of a numeric key by the specified offset.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
// - key: The key to increment.
// - offset: The value to increment by.
//
// Returns:
// - *IncrementResult: The result of the INCREMENT operation.
// - error: Returns an error if the command fails.
func (c *Client) Increment(ctx context.Context, key string, offset int64) (*IncrementResult, error) {
	if c.opts.IsReadonly {
		return nil, fmt.Errorf("cannot execute write op in read-only client: %w", ErrClientReadonly)
	}

	result, err := sendCommand(ctx, c, commandIncr, key, offset)
	if err != nil {
		return nil, err
	}

	if value, ok := result.value.(int64); ok {
		return &IncrementResult{
			NewValue: value,
			Code:     result.code,
		}, nil
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// Decrement decreases the value of a numeric key by the specified offset.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
// - key: The key to decrement.
// - offset: The value to decrement by.
//
// Returns:
// - *DecrementResult: The result of the DECREMENT operation.
// - error: Returns an error if the command fails.
func (c *Client) Decrement(ctx context.Context, key string, offset int64) (*DecrementResult, error) {
	if c.opts.IsReadonly {
		return nil, fmt.Errorf("cannot execute write op in read-only client: %w", ErrClientReadonly)
	}

	result, err := sendCommand(ctx, c, commandDecr, key, offset)
	if err != nil {
		return nil, err
	}

	if value, ok := result.value.(int64); ok {
		return &DecrementResult{
			NewValue: value,
			Code:     result.code,
		}, nil
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// Append adds the specified string to the value of an existing string key.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
// - key: The key whose value will be appended to.
// - value: The string value to append.
//
// Returns:
// - *AppendResult: The result of the APPEND operation.
// - error: Returns an error if the command fails.
func (c *Client) Append(ctx context.Context, key string, value string) (*AppendResult, error) {
	if c.opts.IsReadonly {
		return nil, fmt.Errorf("cannot execute write op in read-only client: %w", ErrClientReadonly)
	}

	result, err := sendCommand(ctx, c, commandAppend, key, value)

	if err != nil {
		return nil, err
	}

	if length, ok := result.value.(int64); ok {
		return &AppendResult{
			ContentLength: length,
			Code:          result.code,
		}, nil
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// MGet retrieves the values of multiple keys from the Universum database.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
// - keys: A slice of keys to retrieve values for.
//
// Returns:
// - *MGetResult: The result of the MGET operation.
// - error: Returns an error if the command fails.
func (c *Client) MGet(ctx context.Context, keys []string) (*MGetResult, error) {
	result, err := sendCommand(ctx, c, commandMget, keys)

	if err != nil {
		return nil, err
	}

	if values, ok := result.value.(map[string]interface{}); ok {
		return &MGetResult{
			Values: values,
			Code:   result.code,
		}, nil
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// MSet sets multiple key-value pairs in the Universum database.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
// - kv: A map of key-value pairs to set in the database.
//
// Returns:
// - *MSetResult: The result of the MSET operation.
// - error: Returns an error if the command fails.
func (c *Client) MSet(ctx context.Context, kv map[string]interface{}) (*MSetResult, error) {
	if c.opts.IsReadonly {
		return nil, fmt.Errorf("cannot execute write op in read-only client: %w", ErrClientReadonly)
	}

	result, err := sendCommand(ctx, c, commandMset, kv)

	if err != nil {
		return nil, err
	}

	if successes, ok := result.value.(map[string]interface{}); ok {
		if converted, err := convertToStringBool(successes); err == nil {
			return &MSetResult{
				Successes: converted,
				Code:      result.code,
			}, nil
		}
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// MDelete deletes multiple keys from the Universum database.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
// - keys: A slice of keys to delete.
//
// Returns:
// - *MDeleteResult: The result of the MDELETE operation.
// - error: Returns an error if the command fails.
func (c *Client) MDelete(ctx context.Context, keys []string) (*MDeleteResult, error) {
	if c.opts.IsReadonly {
		return nil, fmt.Errorf("cannot execute write op in read-only client: %w", ErrClientReadonly)
	}

	result, err := sendCommand(ctx, c, commandMdelete, keys)

	if err != nil {
		return nil, err
	}

	if deletions, ok := result.value.(map[string]interface{}); ok {
		if converted, err := convertToStringBool(deletions); err == nil {
			return &MDeleteResult{
				Deletions: converted,
				Code:      result.code,
			}, nil
		}
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// Info retrieves general information about the Universum database.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
//
// Returns:
// - *InfoResult: The result of the INFO operation.
// - error: Returns an error if the command fails.
func (c *Client) Info(ctx context.Context) (*InfoResult, error) {
	result, err := sendCommand(ctx, c, commandInfo)

	if err != nil {
		return nil, err
	}

	if info, ok := result.value.(string); ok {
		return &InfoResult{
			Raw:  info,
			Code: result.code,
		}, nil
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// Ping sends a ping to the Universum database to check if it is reachable.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
//
// Returns:
// - *PingResult: The result of the PING operation.
// - error: Returns an error if the command fails.
func (c *Client) Ping(ctx context.Context) (*PingResult, error) {
	result, err := sendCommand(ctx, c, commandPing)

	if err != nil {
		return nil, err
	}

	if msg, ok := result.value.(string); ok {
		return &PingResult{
			Message: msg,
			Code:    result.code,
		}, nil
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// TTL retrieves the time-to-live (TTL) value of a specified key.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
// - key: The key to check the TTL for.
//
// Returns:
// - *TTLResult: The result of the TTL operation.
// - error: Returns an error if the command fails.
func (c *Client) TTL(ctx context.Context, key string) (*TTLResult, error) {
	result, err := sendCommand(ctx, c, commandTtl, key)

	if err != nil {
		return nil, err
	}

	if ttl, ok := result.value.(int64); ok {
		return &TTLResult{
			TTL:  time.Duration(ttl) * time.Second,
			Code: result.code,
		}, nil
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// Expire sets the time-to-live (TTL) for a specified key.
//
// Parameters:
// - ctx: Context for managing timeouts and cancellations.
// - key: The key to set the TTL for.
// - ttl: The TTL value in seconds.
//
// Returns:
// - *ExpireResult: The result of the EXPIRE operation.
// - error: Returns an error if the command fails.
func (c *Client) Expire(ctx context.Context, key string, ttl int64) (*ExpireResult, error) {
	if c.opts.IsReadonly {
		return nil, fmt.Errorf("cannot execute write op in read-only client: %w", ErrClientReadonly)
	}

	result, err := sendCommand(ctx, c, commandExpire, key, ttl)

	if err != nil {
		return nil, err
	}

	if success, ok := result.value.(bool); ok {
		return &ExpireResult{
			Success: success,
			Code:    result.code,
		}, nil
	}

	return nil, fmt.Errorf("response value found in expected format: %w", ErrMalformedResponseReceived)
}

// NewClient creates and returns a new Client instance based on the provided options.
// The function initializes the connection pool and generates a unique client ID.
//
// Parameters:
// - opts: A pointer to the Options struct that contains the necessary configurations.
//
// Returns:
// - *Client: A pointer to the newly created Client instance.
// - error: Returns an error if the connection pool could not be initialized.
func NewClient(opts *Options) (*Client, error) {
	ncmu.Lock()
	defer ncmu.Unlock()

	opts.Init()
	 , err := newConnPool(opts)

	if err != nil {
		return nil, err
	}

	currTime := time.Now().UnixNano()
	uniqueId := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(currTime))))

	return &Client{
		id:   uniqueId,
		opts: opts,
		pool: connPool,
	}, nil
}
