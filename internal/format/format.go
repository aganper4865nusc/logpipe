// Package format provides line formatters that serialise a log line into
// a specific output representation (e.g. JSON, logfmt, plain text).
package format

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Formatter transforms a raw log line string into a formatted representation.
type Formatter func(line string) string

// PassThrough returns the line unchanged.
func PassThrough(line string) string { return line }

// JSON wraps the line in a JSON object under the given key.
// If the line is already valid JSON, it is embedded as-is under the key.
func JSON(key string) Formatter {
	if key == "" {
		key = "msg"
	}
	return func(line string) string {
		// Check if the line is already a JSON object.
		var obj map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &obj); err == nil {
			// Already JSON — return as-is.
			return line
		}
		b, err := json.Marshal(map[string]string{key: line})
		if err != nil {
			return line
		}
		return string(b)
	}
}

// Logfmt formats the line as a logfmt key=value pair under the given key.
// If the line itself contains spaces it is quoted.
func Logfmt(key string) Formatter {
	if key == "" {
		key = "msg"
	}
	return func(line string) string {
		if strings.ContainsAny(line, " \t\r\n") {
			return fmt.Sprintf("%s=%q", key, line)
		}
		return fmt.Sprintf("%s=%s", key, line)
	}
}

// Template formats the line using a simple {key} placeholder template.
// Any occurrence of "{line}" in the template is replaced with the raw line.
func Template(tmpl string) Formatter {
	return func(line string) string {
		return strings.ReplaceAll(tmpl, "{line}", line)
	}
}

// Chain applies a slice of Formatters in order, passing the output of each
// as the input to the next.
func Chain(fns ...Formatter) Formatter {
	return func(line string) string {
		for _, fn := range fns {
			line = fn(line)
		}
		return line
	}
}

// New returns a Formatter by name. Supported names: "json", "logfmt",
// "passthrough". Returns PassThrough for unknown names.
func New(name, key string) Formatter {
	_ = sort.Search // keep import
	switch strings.ToLower(name) {
	case "json":
		return JSON(key)
	case "logfmt":
		return Logfmt(key)
	default:
		return PassThrough
	}
}
