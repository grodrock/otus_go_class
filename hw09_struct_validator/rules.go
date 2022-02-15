package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type RuleMatcher interface {
	isMatched(v interface{}) error
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
func (r Rules) isMatched(v interface{}) error {
	for _, rm := range r {
		if err := rm.isMatched(v); err != nil {
			return err
		}
	}
	return nil
}

// GetRuleMatcher return RuleMatcher from rule string
func getRuleMatcher(rulestr string) (RuleMatcher, error) {
	validTagPattern := regexp.MustCompile(`^(.+?):(.*)$`)
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
	case "max":
		if v, err := strconv.Atoi(ruleValue); err == nil {
			return &RuleMaxValidator{v}, nil
		}
		return nil, ErrNotValidRule
	case "regexp":
		rg, err := regexp.Compile(ruleValue)
		if err != nil {
			return nil, RuleRegexpPatternInvalid
		}
		return &RuleRegexpValidator{rg}, nil
	case "in":
		return &RuleInValidator{inVals: ruleValue}, nil

	}

	return nil, ErrNotImplementedRule
}

type RuleLenValidator struct {
	length int
}

func (rvalidator *RuleLenValidator) isMatched(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.String {
		return ErrRuleWrongType
	}
	result := len(v.(string)) == rvalidator.length
	if result {
		return nil
	}
	return RuleLengthInvalid
}

type RuleMinValidator struct {
	min int
}

func (rvalidator *RuleMinValidator) isMatched(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Int {
		return ErrRuleWrongType
	}

	result := v.(int) >= rvalidator.min
	if result {
		return nil
	}
	return RuleMinInvalid
}

type RuleMaxValidator struct {
	max int
}

func (rvalidator *RuleMaxValidator) isMatched(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Int {
		return ErrRuleWrongType
	}

	result := v.(int) <= rvalidator.max
	if result {
		return nil
	}
	return RuleMaxInvalid
}

type RuleRegexpValidator struct {
	pattern *regexp.Regexp
}

func (rvalidator *RuleRegexpValidator) isMatched(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.String {
		return ErrRuleWrongType
	}
	if rvalidator.pattern.MatchString(v.(string)) {
		return nil
	}
	return RuleRegexpInvalid

}

type RuleInValidator struct {
	inVals string
	inMap  map[string]bool
}

func (rvalidator *RuleInValidator) isMatched(v interface{}) error {
	rv := reflect.ValueOf(v)
	// init map
	if rvalidator.inMap == nil {
		rvalidator.inMap = make(map[string]bool)
		for _, k := range strings.Split(rvalidator.inVals, ",") {
			rvalidator.inMap[k] = true
		}
	}
	var key string
	switch rv.Kind() {
	case reflect.String:
		key = v.(string)
	case reflect.Int:
		key = strconv.Itoa(v.(int))
	default:
		return ErrRuleWrongType
	}
	if _, ok := rvalidator.inMap[key]; ok {
		return nil
	}

	return RuleInInvalid
}
