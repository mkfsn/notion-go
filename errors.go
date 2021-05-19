package notion

import (
	"errors"
	"fmt"
)

var (
	ErrUnknown = errors.New("unknown")
)

// ErrorCode https://developers.notion.com/reference/errors
type ErrorCode string

const (
	ErrorCodeInvalidJSON         ErrorCode = "invalid_json"
	ErrorCodeInvalidRequestURI   ErrorCode = "invalid_request_url"
	ErrorCodeInvalidRequest      ErrorCode = "invalid_request"
	ErrorCodeValidationError     ErrorCode = "validation_error"
	ErrorCodeUnauthorized        ErrorCode = "unauthorized"
	ErrorCodeRestrictedResource  ErrorCode = "restricted_resource"
	ErrorCodeObjectNotFound      ErrorCode = "object_not_found"
	ErrorCodeConflictError       ErrorCode = "conflict_error"
	ErrorCodeRateLimited         ErrorCode = "rate_limited"
	ErrorCodeInternalServerError ErrorCode = "internal_server_error"
	ErrorCodeServiceUnavailable  ErrorCode = "service_unavailable"
)

type HTTPError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s", e.Code, e.Message)
}
