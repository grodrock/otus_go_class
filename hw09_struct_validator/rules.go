package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type RuleMatcher interface {
	isMatched(v interface{}) bool
}

type Rules []RuleMatcher

// GetRules get Rules (RuleMatcher) from validation tag
func GetRules(rulesSstr string) (Rules, error) {
	var rules Rules
	for _, r := range strings.Split(rulesSstr, "|") {
		matcher, err := getRuleMatcher(r)
		if err != nil {
			return rules, err
		}
		rules = append(rules, matcher)
	}

	return rules, nil
}

// isMatched validate all rules in Rules is matched to v.
func (r Rules) isMatched(v interface{}) bool {
	for _, rm := range r {
		if ok := rm.isMatched(v); !ok {
			return false
		}
	}
	return true
}

// GetRuleMatcher return RuleMatcher from rule string
func getRuleMatcher(rulestr string) (RuleMatcher, error) {
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
	case "min":
		if v, err := strconv.Atoi(ruleValue); err == nil {
			return &RuleMinValidator{v}, nil
		}
		return nil, ErrNotValidRule

	}

	return nil, ErrNotImplementedRule
}

type RuleLenValidator struct {
	length int
}

func (rvalidator *RuleLenValidator) isMatched(v interface{}) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.String {
		return false
	}
	return len(v.(string)) == rvalidator.length
}

type RuleMinValidator struct {
	min int
}

func (rvalidator *RuleMinValidator) isMatched(v interface{}) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Int {
		return false
	}

	return v.(int) >= rvalidator.min
}
