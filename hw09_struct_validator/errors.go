package hw09structvalidator

import "github.com/pkg/errors"

var (
	// input type validation error.
	ErrNotSupportedType = errors.New("supported types to validate: struct")
	ErrValueValidation  = errors.New("can't validate value")
	// field validation errors.
	ErrFieldValidation = errors.New("field not matches to rules")
	// rule validation errors.
	ErrNotValidRule             = errors.New("not valid rule")
	ErrNotImplementedRule       = errors.New("rule not implemented")
	ErrRuleWrongType            = errors.New("wrong type for this rule")
	ErrRuleLengthInvalid        = errors.New("length rule violation")
	ErrRuleMinInvalid           = errors.New("min value rule violation")
	ErrRuleMaxInvalid           = errors.New("max value rule violation")
	ErrRuleRegexpPatternInvalid = errors.New("regexp pattern invalid")
	ErrRuleRegexpInvalid        = errors.New("regexp value rule violation")
	ErrRuleInInvalid            = errors.New("in value rule violation")
)
