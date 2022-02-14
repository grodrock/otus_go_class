package hw09structvalidator

import "github.com/pkg/errors"

var (
	// input type validation error
	ErrorNotSupportedType = errors.New("supported types to validate: struct")
	// field validation errors
	ErrStringValidation = errors.New("string validation error")
	ErrIntValidation    = errors.New("int validation error")
	// rule validation errors
	ErrNotValidRule       = errors.New("not valid rule")
	ErrNotImplementedRule = errors.New("rule not implemented")
)
