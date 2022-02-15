package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"regexp"
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
			ValidationErrors{ValidationError{"Version", RuleLengthInvalid}},
		},
		{
			in: User{
				ID:     "5e68cffb-6715-4006-8b5e-b50400409005",
				Name:   "User1",
				Age:    25,
				Email:  "example@dot.com",
				Role:   UserRole("admin"),
				Phones: []string{"89001112233"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "5e68cffb-6715-4006-8b5e-b50400409005_xxxx",
				Name:   "User1",
				Age:    99,
				Email:  "exampledot.com",
				Role:   UserRole("janitor"),
				Phones: []string{"89001112233"},
			},
			expectedErr: ValidationErrors{
				ValidationError{"ID", RuleLengthInvalid},
				ValidationError{"Age", RuleMaxInvalid},
				ValidationError{"Email", RuleRegexpInvalid},
				ValidationError{"Role", RuleInInvalid},
			},
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
func TestGetRules(t *testing.T) {
	tests := []struct {
		rulestr       string
		expectedrules Rules
		expectedErr   error
	}{
		{"len:5", Rules{&RuleLenValidator{5}}, nil},
		{"len:5|len:10", Rules{&RuleLenValidator{5}, &RuleLenValidator{10}}, nil},
		{"lenx:5", nil, ErrNotImplementedRule},
		{"len:5x", nil, ErrNotValidRule},
		{"len:5:", nil, ErrNotValidRule},
		{"lenfsdf5", nil, ErrNotValidRule},
		{"min:10", Rules{&RuleMinValidator{10}}, nil},
		{"min:x", nil, ErrNotValidRule},
		{"max:10", Rules{&RuleMaxValidator{10}}, nil},
		{"max:x", nil, ErrNotValidRule},
		{"max:10|len:5|min:20", Rules{&RuleMaxValidator{10}, &RuleLenValidator{5}, &RuleMinValidator{20}}, nil},
		{"regexp:/foo([/", nil, RuleRegexpPatternInvalid},
		{"regexp:^abc&", Rules{&RuleRegexpValidator{regexp.MustCompile(`^abc&`)}}, nil},
		{"in:5,7", Rules{&RuleInValidator{inVals: "5,7"}}, nil},
	}

	for _, tt := range tests {
		matcher, err := GetRules(tt.rulestr)
		require.ErrorIs(t, err, tt.expectedErr)
		require.Equal(t, tt.expectedrules, matcher)
	}
}
func TestIsValid(t *testing.T) {
	tests := []struct {
		in          interface{}
		rulestr     string
		expectedErr error
	}{
		// len
		{"1.2.3", "len:5", nil},
		{"1.2.3", "len:6", RuleLengthInvalid},
		{5, "len:6", ErrRuleWrongType},
		{[]string{"123456", "abcdef"}, "len:6", nil},
		{[]string{"123456", "abcdefx"}, "len:6", RuleLengthInvalid},
		{[]int{1, 2}, "len:6", ErrRuleWrongType},
		{UserRole("admin"), "len:5", nil},
		// min
		{21, "min:20", nil},
		{19, "min:20", RuleMinInvalid},
		{[]int{1, 21, 5}, "min:6", RuleMinInvalid},
		{[]int{12, 21, 28}, "min:6", nil},
		// max
		{21, "max:20", RuleMaxInvalid},
		{19, "max:20", nil},
		{[]int{1, 20, 5}, "max:20", nil},
		{[]int{12, 21, 28}, "max:20", RuleMaxInvalid},
		// min | max
		{19, "min:10|max:20", nil},
		{19, "min:21|max:1", RuleMinInvalid},
		{19, "max:1|min:10", RuleMaxInvalid},
		// regexp
		{"19", "regexp:^Abc", RuleRegexpInvalid},
		{"Abc", "regexp:^Abc", nil},
		{"Abcdef", "regexp:^Abc", nil},
		{"Abcdef", "regexp:^Abc$", RuleRegexpInvalid},
		{"Abcdef", "regexp:^Abc|len:7", RuleLengthInvalid},
		{"Abcdef", "regexp:^Abc|len:6", nil},
		{"some@example.com", "regexp:^\\w+@\\w+\\.\\w+$", nil},
		{UserRole("admin"), "regexp:^a", nil},
		// in
		{"foo", "in:foo", nil},
		{"foo1", "in:foo", RuleInInvalid},
		{[]string{"foo", "bar"}, "in:foo", RuleInInvalid},
		{[]string{"foo", "bar"}, "in:foo,bar", nil},
		{1, "in:1,2,3", nil},
		{[]int{1, 2, 3}, "in:1,2,3,4,5", nil},
		{[]int{1, 2, 3}, "in:1,2,5", RuleInInvalid},
		{UserRole("admin"), "in:admin,stuff", nil},
	}

	for i, tt := range tests {
		rules, err := GetRules(tt.rulestr)
		require.NoError(t, err, "case %d", i)
		isValid := IsValid(tt.in, rules)
		require.Equal(t, tt.expectedErr, isValid,
			"case %d validation miss: %s against %s", i, tt.in, tt.rulestr)
	}
}
