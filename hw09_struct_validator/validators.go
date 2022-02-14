package hw09structvalidator

import (
	"reflect"
)

func IsValid(v interface{}, rm RuleMatcher) bool {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return rm.isMatched(rv.String())
	}

	return rm.isMatched(v)

	// return false
}
