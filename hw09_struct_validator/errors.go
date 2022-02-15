package hw09structvalidator

import "github.com/pkg/errors"

var (
	// input type validation error
	ErrNotSupportedType = errors.New("supported types to validate: struct")
	ErrValueValidation  = errors.New("can't validate value")
	// field validation errors
	ErrFieldValidation = errors.New("field not matches to rules")
	// rule validation errors
	ErrNotValidRule       = errors.New("not valid rule")
	ErrNotImplementedRule = errors.New("rule not implemented")
	ErrRuleWrongType      = errors.New("wrong type for this rule")
	RuleLengthInvalid     = errors.New("length rule violation")
	RuleMinInvalid        = errors.New("min value rule violation")
	RuleMaxInvalid        = errors.New("max value rule violation")
)
