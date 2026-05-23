// Package schema provides JSON log line validation against a required field set.
package schema

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Validator checks that a log line satisfies a required set of JSON fields.
type Validator struct {
	required []string
}

// ValidationError is returned when one or more required fields are absent.
type ValidationError struct {
	Missing []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("schema: missing required fields: %s", strings.Join(e.Missing, ", "))
}

// New returns a Validator that enforces the given required field names.
// An empty required list creates a pass-through validator.
func New(required []string) *Validator {
	copy := make([]string, len(required))
	copy = append(copy[:0], required...)
	return &Validator{required: copy}
}

// Validate checks that line is valid JSON and contains every required field.
// Non-JSON lines are rejected when at least one required field is configured.
func (v *Validator) Validate(line string) error {
	if len(v.required) == 0 {
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return fmt.Errorf("schema: line is not valid JSON: %w", err)
	}

	var missing []string
	for _, field := range v.required {
		if _, ok := obj[field]; !ok {
			missing = append(missing, field)
		}
	}

	if len(missing) > 0 {
		return &ValidationError{Missing: missing}
	}
	return nil
}

// PassThrough returns a Validator that accepts every line without inspection.
func PassThrough() *Validator {
	return &Validator{}
}
