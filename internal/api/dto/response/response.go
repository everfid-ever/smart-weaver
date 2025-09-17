package response

import "encoding/json"

// Response defines a common response structure
type Response[T any] struct {
	Code string `json:"code"`
	Info string `json:"info"`
	Data T      `json:"data"`
}

// NewResponse creates a response object
func NewResponse[T any](code, info string, data T) *Response[T] {
	return &Response[T]{
		Code: code,
		Info: info,
		Data: data,
	}
}

// Success returns success response
func Success[T any](data T) *Response[T] {
	return NewResponse("0000", "成功", data)
}

// Error returns error response
func Error[T any](code, info string) *Response[T] {
	var zero T
	return NewResponse(code, info, zero)
}

// String returns JSON string
func (r *Response[T]) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}
