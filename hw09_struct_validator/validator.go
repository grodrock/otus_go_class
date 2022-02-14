package hw09structvalidator

import (
	"log"
	"reflect"

	"github.com/pkg/errors"
)

type ValidationError struct {
	Field string
	Err   error
}

const validateTag string = "validate"

// errors
var ()

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
		return ErrorNotSupportedType
	}
	return ValidateStruct(rv)
}

func ValidateStruct(rv reflect.Value) error {
	var validationErrors ValidationErrors
	t := rv.Type()
	fieldsNum := t.NumField()
	log.Printf("type: %v, fields: %d", t, fieldsNum)

	for i := 0; i < fieldsNum; i++ {
		field := t.Field(i)
		fv := rv.Field(i)
		fieldName := field.Name

		// get validation rule if exist
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
		// validate value
		if !IsValid(fv, rules) {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   ErrFieldValidation,
			})
		}

	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
