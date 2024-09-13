package universum

import (
	"bufio"

	"github.com/cshekharsharma/resp-go/resp3"
)

// Wrapper function for resp3.Encode
func encodeResp(value interface{}) (string, error) {
	return resp3.Encode(value)
}

// Wrapper function for resp3.Decode
func decodeResp(reader *bufio.Reader) (interface{}, error) {
	return resp3.Decode(reader)
}
