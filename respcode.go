package universum

// Universum server response codes
const (
	RespPingSuccess     int64 = 200
	RespSnapshotStarted int64 = 201

	RespServerShuttingDown int64 = 501
	RespServerBusy         int64 = 502

	RespRecordFound   int64 = 1000
	RespRecordUpdated int64 = 1001
	RespRecordDeleted int64 = 1002
	RespHelpContentOk int64 = 1010
	RespInfoContentOk int64 = 1011

	RespMgetCompleted int64 = 1100
	RespMsetCompleted int64 = 1101
	RespMdelCompleted int64 = 1102

	RespInvalidCmdInput  int64 = 5000
	RespRecordNotFound   int64 = 5001
	RespRecordExpired    int64 = 5002
	RespRecordNotDeleted int64 = 5003
	RespIncrInvalidType  int64 = 5004

	RespRecordTooBig     int64 = 5005
	RespIinvalidDatatype int64 = 5006
)
