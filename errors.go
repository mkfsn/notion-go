package notion

import (
	"errors"
	"fmt"
)

var (
	ErrUnknown = errors.New("unknown")
)

type HTTPError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s", e.Code, e.Message)
}
