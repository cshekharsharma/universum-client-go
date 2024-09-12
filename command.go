package universum

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"time"
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

const remoteByteDelimiter = "\x04\x04\x04\x04"

func sendCommand(ctx context.Context, c *Client, command string, args ...interface{}) (*CommandResult, error) {
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

	if c.opts.WriteTimeout > 0 {
		err := conn.getNetConn().SetWriteDeadline(time.Now().Add(c.opts.WriteTimeout))
		if err != nil {
			return nil, fmt.Errorf("failed to set write deadline: %v", err)
		}
	}

	bytesWritten, err := conn.write([]byte(encodedCommand))
	if err != nil {
		return nil, fmt.Errorf("failed while writing bytes to the socket: %w", ErrSocketWriteFailed)
	}
	if bytesWritten != len(encodedCommand) {
		return nil, fmt.Errorf("incomplete write: wrote %d/%d bytes: %w",
			bytesWritten, len(encodedCommand), ErrIncompleteSocketWrite)
	}

	if err := conn.getWriter().Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush writer: %w", ErrSocketFlushFailed)
	}

	decodedBuffer, err := readUntilDelimiter(conn, c.opts, remoteByteDelimiter)
	if err != nil {
		return nil, fmt.Errorf("failed while reading bytes from the socket: [%v] %w", err, ErrSocketReadFailed)
	}

	decoded, err := decodeResp(bufio.NewReader(decodedBuffer))
	if _, ok := decoded.(error); ok {
		return nil, fmt.Errorf("server rejected the request: %v : %w", decoded, ErrServerRejectedRequested)
	}

	if err != nil {
		return nil, err
	}

	return toCommandResult(decoded)
}

func readUntilDelimiter(conn connInterface, opts *Options, delim string) (*bytes.Buffer, error) {
	delimiterBytes := []byte(delim)
	delimiterLen := len(delimiterBytes)

	if opts.ReadTimeout > 0 {
		err := conn.getNetConn().SetReadDeadline(time.Now().Add(opts.ReadTimeout))
		if err != nil {
			return nil, fmt.Errorf("failed to set read deadline: %v", err)
		}
	}

	reader := conn.getReader()
	var buffer bytes.Buffer
	var decoderBuffer bytes.Buffer
	pipe := make([]byte, 1024) // Read in chunks

	for {
		chunkSize, err := reader.Read(pipe)
		if err != nil {
			if err == io.EOF && buffer.Len() > 0 {
				// If EOF is received but some data is present, return what we have.
				break
			}
			return nil, fmt.Errorf("error reading from connection: %v", err)
		}

		buffer.Write(pipe[:chunkSize])
		decoderBuffer.Write(pipe[:chunkSize])

		if buffer.Len() >= delimiterLen && bytes.HasSuffix(buffer.Bytes(), delimiterBytes) {
			decoderBuffer.Truncate(decoderBuffer.Len() - delimiterLen)
			break
		}
	}

	return &decoderBuffer, nil
}
