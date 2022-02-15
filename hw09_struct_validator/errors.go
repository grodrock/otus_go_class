package hw09structvalidator

import "github.com/pkg/errors"

var (
	// input type validation error
	ErrNotSupportedType = errors.New("supported types to validate: struct")
	// field validation errors
	ErrStringValidation = errors.New("string validation error")
	ErrIntValidation    = errors.New("int validation error")
	ErrFieldValidation  = errors.New("field not matches to rules")
	// rule validation errors
	ErrNotValidRule       = errors.New("not valid rule")
	ErrNotImplementedRule = errors.New("rule not implemented")
)
