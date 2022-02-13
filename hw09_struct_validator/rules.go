package hw09structvalidator

import (
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

var (
	ErrNotValidRule       = errors.New("not valid rule")
	ErrNotImplementedRule = errors.New("rule not implemented")
)

type RuleMatcher interface {
	isMatched(v interface{}) bool
}

type RuleLenValidator struct {
	length int
}

func (rv *RuleLenValidator) isMatched(v interface{}) bool {
	val, ok := v.(string)
	if !ok {
		return false
	}
	return len(val) == rv.length
}

func GetRuleMatcher(rulestr string) (RuleMatcher, error) {
	validTagPattern, _ := regexp.Compile(`(.*):(.*)`)
	if !validTagPattern.MatchString(rulestr) {
		return nil, ErrNotValidRule
	}

	fields := validTagPattern.FindStringSubmatch(rulestr)
	ruleName := fields[1]
	ruleValue := fields[2]
	switch ruleName {
	case "len":
		if v, err := strconv.Atoi(ruleValue); err == nil {
			return &RuleLenValidator{v}, nil
		}
		return nil, ErrNotValidRule

	}

	return nil, ErrNotImplementedRule
}
