package hw09structvalidator

import (
	"reflect"
)

// wrapper on rules validation: helps to validate slice elems
func IsValid(v interface{}, rm RuleMatcher) error {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice:

		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i)
			if !elem.CanInterface() {
				return ErrValueValidation
			}
			err := rm.isMatched(elem.Interface())
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return rm.isMatched(v)
	}

	// return false
}
