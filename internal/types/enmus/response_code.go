package enmus

// ResponseCodeEnum defines the business response code
type ResponseCodeEnum struct {
	Code string
	Info string
}

var (
	EnumSuccess      = ResponseCodeEnum{"0000", "Success"}
	EnumUnError      = ResponseCodeEnum{"0001", "Unknown failure"}
	EnumIllegalParam = ResponseCodeEnum{"0002", "Illegal parameters"}
)
