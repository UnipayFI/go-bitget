package client

import "fmt"

// APIError is the error envelope Bitget returns when a request fails. Bitget
// codes are strings; "00000" means success and anything else is an error.
type APIError struct {
	Code        string `json:"code"`
	Message     string `json:"msg"`
	RequestTime int64  `json:"requestTime"`
}

// Error returns the error code and message.
func (e APIError) Error() string {
	return fmt.Sprintf("<APIError> code=%s, msg=%s", e.Code, e.Message)
}

// IsValid reports whether e represents an actual API-level error (a non-empty,
// non-success code).
func (e APIError) IsValid() bool {
	return e.Code != "" && e.Code != "00000"
}

// IsAPIError reports whether err is a Bitget *APIError.
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}
