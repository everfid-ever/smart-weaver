package common

const (
	SPLIT = ","
)

// ResponseCode defines the response code
type ResponseCode struct {
	Code string
	Info string
}

var (
	ResponseSuccess        = ResponseCode{"0000", "Success"}
	ResponseUnError        = ResponseCode{"0001", "Unknown failure"}
	ResponseIllegalParam   = ResponseCode{"0002", "Illegal parameters"}
	ResponseIndexException = ResponseCode{"0003", "Unique index conflict"}
	ResponseUpdateZero     = ResponseCode{"0004", "Update record is 0"}
	ResponseHttpException  = ResponseCode{"0005", "HTTP interface call exception"}
)
