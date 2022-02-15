package hw09structvalidator

import (
	"reflect"

	"github.com/pkg/errors"
)

type ValidationError struct {
	Field string
	Err   error
}

const validateTag string = "validate"

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var wrappedErr error

	for _, ve := range v {
		wrappedErr = errors.Wrap(wrappedErr, ve.Field)
		// wrappedErr = fmt.Errorf(ve.Field, ve.Err)
	}

	return wrappedErr.Error()
}

func Validate(v interface{}) error {
	// check value is a struct
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return ErrNotSupportedType
	}
	return ValidateStruct(rv)
}

func ValidateStruct(rv reflect.Value) error {
	var validationErrors ValidationErrors
	t := rv.Type()
	fieldsNum := t.NumField()

	for i := 0; i < fieldsNum; i++ {
		field := t.Field(i)
		fv := rv.Field(i)
		fieldName := field.Name

		// get validation rule string if exist
		validateTagString, ok := field.Tag.Lookup(validateTag)
		if !ok {
			continue
		}

		// create rules from validateTagString
		rules, err := GetRules(validateTagString)
		if err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   err,
			})
			continue
		}

		// check if we can wrap it to interface{}
		if !fv.CanInterface() {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   ErrValueValidation,
			})
		}

		// validate value against rule
		err = IsValid(fv.Interface(), rules)
		if err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
