package hw09structvalidator

import (
	"fmt"

	"github.com/pkg/errors"
)

func ValidateString(value string, rule string) error {
	fmt.Println("Validating string:", value, "with rule:", rule)

	ruleMatcher, err := GetRuleMatcher(rule)
	if err != nil {
		return err
	}
	if ruleMatcher == nil {
		// test case - we should always have RuleMatcher here
		return errors.New("nor error or RuleMatcher recieved from GetRuleMatcher")
	}
	if !ruleMatcher.isMatched(value) {
		return ErrStringValidation
	}

	return nil
}
