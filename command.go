package universum

import (
	"fmt"
)

const (
	commandPing     string = "PING"
	commandExists   string = "EXISTS"
	commandGet      string = "GET"
	commandSet      string = "SET"
	commandDelete   string = "DELETE"
	commandIncr     string = "INCR"
	commandDecr     string = "DECR"
	commandAppend   string = "APPEND"
	commandMget     string = "MGET"
	commandMset     string = "MSET"
	commandMdelete  string = "MDELETE"
	commandTtl      string = "TTL"
	commandExpire   string = "EXPIRE"
	commandSnapshot string = "SNAPSHOT"
	commandInfo     string = "INFO"
	commandHelp     string = "HELP"
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
