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

// Universum error codes
const (
	RespPingSuccess     uint32 = 200
	RespSnapshotStarted uint32 = 201

	RespServerShuttingDown uint32 = 501
	RespServerBusy         uint32 = 502

	RespRecordFound   uint32 = 1000
	RespRecordUpdated uint32 = 1001
	RespRecordDeleted uint32 = 1002
	RespHelpContentOk uint32 = 1010
	RespInfoContentOk uint32 = 1011

	RespMgetCompleted uint32 = 1100
	RespMsetCompleted uint32 = 1101
	RespMdelCompleted uint32 = 1102

	RespInvalidCmdInput  uint32 = 5000
	RespRecordNotFound   uint32 = 5001
	RespRecordExpired    uint32 = 5002
	RespRecordNotDeleted uint32 = 5003
	RespIncrInvalidType  uint32 = 5004

	RespRecordTooBig     uint32 = 5005
	RespIinvalidDatatype uint32 = 5006
)

var (
	errUnexpectedRead = errors.New("unexpected read from socket")
)

func isWriteableDatatype(value interface{}) bool {
	switch value.(type) {
	case string,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, bool:
		return true

	case []interface{}:
		return true

	default:
		return false
	}
}
