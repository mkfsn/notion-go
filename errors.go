package notion

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrUnimplemented = errors.New("unimplemented")
	ErrUnknown       = errors.New("unknown")
)

type HTTPError struct {
	StatusCode int
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("StatusCode: %d, Code: %s, Message: %s", e.StatusCode, e.Code, e.Message)
}

func newHTTPError(statusCode int, data []byte) error {
	var httpError HTTPError

	if err := json.Unmarshal(data, &httpError); err != nil {
		return ErrUnknown
	}

	httpError.StatusCode = statusCode

	return httpError
}
