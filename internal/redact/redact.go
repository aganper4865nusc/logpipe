// Package redact provides log line redaction for sensitive data patterns
// such as passwords, tokens, and credit card numbers.
package redact

import (
	"fmt"
	"regexp"
)

// Rule defines a named redaction rule with a compiled regex pattern
// and a replacement string.
type Rule struct {
	Name        string
	Pattern     *regexp.Regexp
	Replacement string
}

// Redactor applies a set of redaction rules to log lines.
type Redactor struct {
	rules []Rule
}

// New creates a Redactor from the provided pattern map.
// Keys are rule names; values are regex patterns.
// Returns an error if any pattern fails to compile.
func New(patterns map[string]string) (*Redactor, error) {
	rules := make([]Rule, 0, len(patterns))
	for name, pat := range patterns {
		re, err := regexp.Compile(pat)
		if err != nil {
			return nil, fmt.Errorf("redact: invalid pattern %q: %w", name, err)
		}
		rules = append(rules, Rule{
			Name:        name,
			Pattern:     re,
			Replacement: "[REDACTED]",
		})
	}
	return &Redactor{rules: rules}, nil
}

// Transform applies all redaction rules to line and returns the sanitised result.
func (r *Redactor) Transform(line string) string {
	for _, rule := range r.rules {
		line = rule.Pattern.ReplaceAllString(line, rule.Replacement)
	}
	return line
}

// PassThrough returns the line unchanged. Useful as a no-op when redaction
// is disabled.
func PassThrough(line string) string { return line }
