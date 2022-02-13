package hw09structvalidator

import (
	"reflect"

	"github.com/pkg/errors"
)

type ValidationError struct {
	Field string
	Err   error
}

var (
	ErrorNotStructure     error = errors.New("not a structure")
	ErrorStringValidation error = errors.New("string validation error")
	ErrorIntValidation    error = errors.New("int validation error")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var wrappedErr error

	for _, ve := range v {
		wrappedErr = errors.Wrap(wrappedErr, ve.Field)
		//wrappedErr = fmt.Errorf(ve.Field, ve.Err)
	}

	return wrappedErr.Error()
}

func Validate(v interface{}) error {
	// check value is a struct
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return ErrorNotStructure
	}
	return nil
}
