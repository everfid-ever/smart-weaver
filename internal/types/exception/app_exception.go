package types

import "fmt"

// AppException customs application exception
type AppException struct {
	Code string
	Info string
	Err  error
}

func NewAppException(code string) *AppException {
	return &AppException{Code: code}
}

func NewAppExceptionWithCause(code string, cause error) *AppException {
	return &AppException{Code: code, Err: cause}
}

func NewAppExceptionWithMessage(code, message string) *AppException {
	return &AppException{Code: code, Info: message}
}

func NewAppExceptionFull(code, message string, cause error) *AppException {
	return &AppException{Code: code, Info: message, Err: cause}
}

func (e *AppException) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("AppException{code='%s', info='%s', cause='%v'}", e.Code, e.Info, e.Err)
	}
	return fmt.Sprintf("AppException{code='%s', info='%s'}", e.Code, e.Info)
}
