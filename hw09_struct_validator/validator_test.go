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
		{UserRole("role"), ErrNotSupportedType},
		{App{"1.2.3"}, nil},
		{
			App{"1.2.34"},
			ValidationErrors{ValidationError{"Version", ErrFieldValidation}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if err != nil {
				require.Equal(t, tt.expectedErr, err, "err %v not as expected %v", err, tt.expectedErr)
			} else {
				require.Nil(t, tt.expectedErr)
			}

		})
	}
}
func TestRuleMatcher(t *testing.T) {
	tests := []struct {
		rulestr       string
		expectedrules Rules
		expectedErr   error
	}{
		{"len:5", Rules{&RuleLenValidator{5}}, nil},
		{"len:5|len:10", Rules{&RuleLenValidator{5}, &RuleLenValidator{10}}, nil},
		{"lenx:5", nil, ErrNotImplementedRule},
		{"len:5x", nil, ErrNotValidRule},
		{"lenfsdf5", nil, ErrNotValidRule},
		{"min:10", Rules{&RuleMinValidator{10}}, nil},
		{"min:x", nil, ErrNotValidRule},
	}

	for _, tt := range tests {
		matcher, err := GetRules(tt.rulestr)
		require.ErrorIs(t, err, tt.expectedErr)
		require.Equal(t, tt.expectedrules, matcher)
	}
}
func TestIsValid(t *testing.T) {
	tests := []struct {
		in              interface{}
		rulestr         string
		expectedisValid bool
	}{
		{"1.2.3", "len:5", true},
		{"1.2.3", "len:6", false},
		{5, "len:6", false},
		{21, "min:20", true},
		{19, "min:20", false},
	}

	for _, tt := range tests {
		rules, err := GetRules(tt.rulestr)
		require.NoError(t, err)
		isValid := IsValid(tt.in, rules)
		require.Equal(t, tt.expectedisValid, isValid,
			"validation miss: %s against %s", tt.in, tt.rulestr)
	}
}
