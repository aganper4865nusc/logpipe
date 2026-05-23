// Package parser provides log line parsers that convert raw strings into
// structured key-value maps for downstream processing.
package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ParseFn transforms a raw log line into a structured map.
// Returns an error if the line cannot be parsed according to the format.
type ParseFn func(line string) (map[string]any, error)

// JSON returns a ParseFn that unmarshals JSON log lines.
func JSON() ParseFn {
	return func(line string) (map[string]any, error) {
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			return nil, fmt.Errorf("parser: invalid JSON: %w", err)
		}
		return m, nil
	}
}

// Regex returns a ParseFn that extracts named capture groups from a line.
// Each named group becomes a key in the resulting map.
func Regex(pattern string) (ParseFn, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("parser: invalid pattern: %w", err)
	}
	names := re.SubexpNames()
	return func(line string) (map[string]any, error) {
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			return nil, fmt.Errorf("parser: line did not match pattern")
		}
		m := make(map[string]any, len(names))
		for i, name := range names {
			if i == 0 || name == "" {
				continue
			}
			m[name] = matches[i]
		}
		return m, nil
	}, nil
}

// KV returns a ParseFn that parses key=value pairs separated by pairSep.
// Pairs themselves are split on kvSep (e.g. "=" or ":").
func KV(pairSep, kvSep string) ParseFn {
	return func(line string) (map[string]any, error) {
		m := make(map[string]any)
		for _, pair := range strings.Split(line, pairSep) {
			parts := strings.SplitN(strings.TrimSpace(pair), kvSep, 2)
			if len(parts) != 2 {
				continue
			}
			m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
		if len(m) == 0 {
			return nil, fmt.Errorf("parser: no key-value pairs found")
		}
		return m, nil
	}
}

// PassThrough returns a ParseFn that wraps the raw line under the key "message".
func PassThrough() ParseFn {
	return func(line string) (map[string]any, error) {
		return map[string]any{"message": line}, nil
	}
}
