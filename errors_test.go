package notion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPError_Error(t *testing.T) {
	type fields struct {
		Code    ErrorCode
		Message string
	}
	tests := []struct {
		fields    fields
		wantError string
	}{
		{
			fields: fields{
				Code:    ErrorCodeInvalidJSON,
				Message: "missing ]",
			},
			wantError: "Code: invalid_json, Message: missing ]",
		},
	}
	for _, tt := range tests {
		t.Run(string(tt.fields.Code), func(t *testing.T) {
			e := HTTPError{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
			}

			assert.EqualError(t, e, tt.wantError)
		})
	}
}
