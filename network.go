package universum

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"
)

const remoteByteDelimiter = "\x04\x04\x04\x04"

func sendCommand(ctx context.Context, c *Client, command string, args ...interface{}) (interface{}, error) {
	conn, err := c.pool.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	defer c.pool.ReleaseConn(ctx, conn)

	cmdInput := make([]interface{}, 0, len(args)+1)
	cmdInput = append(cmdInput, command)
	cmdInput = append(cmdInput, args...)

	encodedCommand, err := encodeResp(cmdInput)
	if err != nil {
		return nil, fmt.Errorf("resp encoding failed before sending the command: %w", ErrCommandEncodingFailed)
	}

	if _, err := conn.write([]byte(encodedCommand)); err != nil {
		return nil, fmt.Errorf("failed while writing bytes to the socket: %w", ErrSocketWriteFailed)
	}

	err = readUntilDelimiter(conn, c.opts, remoteByteDelimiter)
	if err != nil {
		return nil, fmt.Errorf("failed while reading bytes from the socket: [%v] %w", err, ErrSocketReadFailed)
	}

	return decodeResp(conn.getReader())
}

func readUntilDelimiter(conn connInterface, opts *Options, delim string) error {
	delimiterBytes := []byte(delim)
	delimiterLen := len(delimiterBytes)

	var buffer bytes.Buffer
	pipe := make([]byte, 1024) // Read in chunks

	if opts.ReadTimeout > 0 {
		err := conn.getNetConn().SetReadDeadline(time.Now().Add(opts.ReadTimeout))
		if err != nil {
			return fmt.Errorf("failed to set read deadline: %v", err)
		}
	}

	for {
		reader := conn.getReader()
		chunk, err := reader.Read(pipe)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error reading from connection: %v", err)
		}

		buffer.Write(pipe[:chunk])

		if buffer.Len() >= delimiterLen {
			if bytes.Equal(buffer.Bytes()[buffer.Len()-delimiterLen:], delimiterBytes) {
				// We've found the delimiter, and we stop here
				return nil
			}
		}
	}
}
