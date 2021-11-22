package collectors

import (
	"regexp"
	"strings"
)

type Filter struct {
	rules []*regexp.Regexp
}

func NewFilter(rules []string) (*Filter, error) {
	filter := &Filter{
		rules: []*regexp.Regexp{},
	}
	for _, rule := range rules {
		compiledRegexp, err := regexp.Compile(rule)
		if err != nil { return nil, err }
		filter.rules = append(filter.rules, compiledRegexp)
	}
	return filter, nil
}

func (f *Filter) Filter(name string) bool {
	for _, rule := range f.rules {
		match := rule.MatchString(strings.TrimSpace(name))
		if match {
			return true
		}
	}
	return false
}

func (f *Filter) Size() int {
	return len(f.rules)
}

