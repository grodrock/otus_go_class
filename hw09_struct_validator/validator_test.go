package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{UserRole("role"), ErrorNotSupportedType},
		{App{"1.2.3"}, nil},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr, "err %v not as expected %v", err, tt.expectedErr)

		})
	}
}

func TestStringValidator(t *testing.T) {
	tests := []struct {
		in          string
		rule        string
		expectedErr error
	}{
		{"1.2.3", "len:5", nil},
		{"1.2.3", "len:6", ErrStringValidation},
		{"1.2.34", "len:5", ErrStringValidation},
		{"1.2.3", "len:x", ErrNotValidRule},
	}

	for _, tt := range tests {
		err := ValidateString(tt.in, tt.rule)
		require.ErrorIs(t, err, tt.expectedErr)
	}
}
