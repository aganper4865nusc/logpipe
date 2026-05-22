package filter

import (
	"regexp"
	"strings"
)

// Rule defines a single filtering rule applied to log lines.
type Rule struct {
	Contains string `yaml:"contains,omitempty"`
	Regex    string `yaml:"regex,omitempty"`
	Level    string `yaml:"level,omitempty"`
}

// Filter holds compiled rules and decides whether a log line should pass.
type Filter struct {
	rules    []Rule
	compiled []*regexp.Regexp
}

// New creates a Filter from the given rules, pre-compiling any regex patterns.
func New(rules []Rule) (*Filter, error) {
	f := &Filter{rules: rules}
	for _, r := range rules {
		if r.Regex != "" {
			re, err := regexp.Compile(r.Regex)
			if err != nil {
				return nil, err
			}
			f.compiled = append(f.compiled, re)
		} else {
			f.compiled = append(f.compiled, nil)
		}
	}
	return f, nil
}

// Match returns true if the line satisfies ALL defined rules.
func (f *Filter) Match(line string) bool {
	for i, rule := range f.rules {
		if rule.Contains != "" && !strings.Contains(line, rule.Contains) {
			return false
		}
		if rule.Level != "" && !strings.Contains(strings.ToLower(line), strings.ToLower(rule.Level)) {
			return false
		}
		if f.compiled[i] != nil && !f.compiled[i].MatchString(line) {
			return false
		}
	}
	return true
}

// PassThrough returns a Filter that matches every line.
func PassThrough() *Filter {
	return &Filter{}
}
