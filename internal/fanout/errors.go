package fanout

import (
	"fmt"
	"strings"
)

// MultiError aggregates one or more errors returned by fanout sinks.
type MultiError struct {
	Errs []error
}

// Error implements the error interface, joining all messages with a semicolon.
func (m *MultiError) Error() string {
	msgs := make([]string, 0, len(m.Errs))
	for _, e := range m.Errs {
		msgs = append(msgs, e.Error())
	}
	return fmt.Sprintf("fanout: %d sink error(s): %s", len(m.Errs), strings.Join(msgs, "; "))
}

// Unwrap returns the slice of wrapped errors so errors.Is / errors.As work
// across the full list.
func (m *MultiError) Unwrap() []error {
	return m.Errs
}
