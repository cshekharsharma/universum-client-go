package universum

import "errors"

// Networking errors
var (
	ErrConnectionDialFailed   = errors.New("CONN_DIAL_FAILED")
	ErrConnectionDialTimeout  = errors.New("CONN_DIAL_TIMEOUT")
	ErrConnectionWaitTimeout  = errors.New("CONN_WAIT_TIMEOUT")
	ErrConnectionConfigFailed = errors.New("CONN_CONFIG_FAILED")

	ErrConnectionPoolClosed = errors.New("CONN_POOL_CLOSED")

	ErrCommandEncodingFailed = errors.New("CMD_ENCODING_FAILED")
	ErrSocketWriteFailed     = errors.New("SOCKET_WRITE_FAILED")
	ErrIncompleteSocketWrite = errors.New("INCOMPLETE_SOCKET_WRITE")
	ErrSocketFlushFailed     = errors.New("SOCKET_FLUSH_FAILED")
	ErrSocketReadFailed      = errors.New("SOCKET_READ_FAILED")

	ErrMalformedResponseReceived = errors.New("MALFORMED_RESPONSE_RECEIVED")
	ErrServerRejectedRequest     = errors.New("SERVER_REJECTED_REQUEST")

	ErrInvalidRequest  = errors.New("INVALID_REQUEST")
	ErrClientReadonly  = errors.New("CLIENT_READONLY")
	ErrInvalidDatatype = errors.New("INVALID_DATATYPE")
)

var (
	errUnexpectedRead = errors.New("unexpected read from socket")
)
