package hw09structvalidator

import (
	"fmt"
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
var (
	ErrorNotStructure error = errors.New("not a structure")
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
	return ValidateStruct(rv)
}

func ValidateStruct(rv reflect.Value) error {
	t := rv.Type()
	fieldsNum := t.NumField()
	log.Printf("type: %v, fields: %d", t, fieldsNum)

	for i := 0; i < fieldsNum; i++ {
		field := t.Field(i)
		fv := rv.Field(i)

		// get validation rule if exist
		rulesString, ok := field.Tag.Lookup(validateTag)
		if !ok {
			continue
		}

		switch fv.Kind() {
		case reflect.String:
			err := ValidateString(fv.String(), rulesString)
			fmt.Println(err)
		}

	}

	return nil
}
