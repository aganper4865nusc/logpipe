// Package masker provides field-level masking for structured log lines.
// It replaces the value of named JSON keys with a configurable mask string,
// allowing sensitive fields (e.g. passwords, tokens) to be hidden before
// forwarding to a sink.
package masker

import (
	"encoding/json"
	"strings"
)

const defaultMask = "***"

// Masker replaces values of configured field names in JSON log lines.
type Masker struct {
	fields map[string]struct{}
	mask   string
}

// New returns a Masker that will replace values for the given field names.
// mask is the replacement string; if empty, defaultMask ("***") is used.
func New(fields []string, mask string) *Masker {
	if mask == "" {
		mask = defaultMask
	}
	fm := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		fm[strings.TrimSpace(f)] = struct{}{}
	}
	return &Masker{fields: fm, mask: mask}
}

// Transform masks configured fields in a JSON log line.
// Non-JSON lines are returned unchanged.
func (m *Masker) Transform(line string) string {
	if len(m.fields) == 0 {
		return line
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for k := range obj {
		if _, ok := m.fields[k]; ok {
			obj[k] = m.mask
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// PassThrough returns a transform function that leaves lines unchanged.
// Useful as a no-op placeholder when masking is disabled.
func PassThrough(line string) string { return line }
