package hw09structvalidator

import (
	"reflect"
)

// wrapper on rules validation: helps to validate slice elems
func IsValid(v interface{}, rm RuleMatcher) bool {
	rv := reflect.ValueOf(v)
	var isMatched bool
	switch rv.Kind() {
	case reflect.Slice:

		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i)
			if !elem.CanInterface() {
				return false
			}
			isMatched = rm.isMatched(elem.Interface())
			if !isMatched {
				return false
			}
		}
		return true
	default:
		return rm.isMatched(v)
	}

	// return false
}
