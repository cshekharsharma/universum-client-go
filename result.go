package universum

import (
	"fmt"
	"time"
)

type CommandResult struct {
	value   interface{}
	code    int64
	message string
}

func toCommandResult(result interface{}) (*CommandResult, error) {
	if result == nil {
		return nil, fmt.Errorf("empty result from server found: %w", ErrMalformedResponseReceived)
	}

	if res, ok := result.([]interface{}); ok && len(res) == 3 {
		value := res[0]
		code, codeOk := res[1].(int64)
		message, msgOk := res[2].(string)

		if codeOk && msgOk {
			return &CommandResult{
				value:   value,
				code:    code,
				message: message,
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid result from server found: %w", ErrMalformedResponseReceived)
}

type GetResult struct {
	Value interface{}
	Code  int64
}

type SetResult struct {
	Success bool
	Code    int64
}

type ExistsResult struct {
	Found bool
	Code  int64
}

type DeleteResult struct {
	Deleted bool
	Code    int64
}

type IncrementResult struct {
	NewValue int64
	Code     int64
}

type DecrementResult struct {
	NewValue int64
	Code     int64
}

type AppendResult struct {
	ContentLength int64
	Code          int64
}

type MGetResult struct {
	Values map[string]interface{}
	Code   int64
}

type MSetResult struct {
	Successes map[string]bool
	Code      int64
}

type MDeleteResult struct {
	Deletions map[string]bool
	Code      int64
}

type ExpireResult struct {
	Success bool
	Code    int64
}

type TTLResult struct {
	TTL  time.Duration
	Code int64
}

type PingResult struct {
	Message string
	Code    int64
}

type InfoResult struct {
	Raw  string
	Code int64
}
