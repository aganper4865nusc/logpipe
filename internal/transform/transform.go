// Package transform provides line transformation functions that can be
// applied to log entries before they are forwarded to sinks.
package transform

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Transformer mutates a log line and returns the result.
// If the transformer returns an error the line is dropped.
type Transformer func(line string) (string, error)

// Chain combines multiple transformers into one. Each transformer is applied
// in order; if any returns an error the chain stops and the error is returned.
func Chain(transformers ...Transformer) Transformer {
	return func(line string) (string, error) {
		var err error
		for _, t := range transformers {
			line, err = t(line)
			if err != nil {
				return "", err
			}
		}
		return line, nil
	}
}

// AddTimestamp wraps the line in a JSON envelope with a UTC timestamp field.
// If the line is already valid JSON the timestamp key is injected into the
// existing object; otherwise the line is stored under the "message" key.
func AddTimestamp(key string) Transformer {
	return func(line string) (string, error) {
		ts := time.Now().UTC().Format(time.RFC3339Nano)
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			obj = map[string]interface{}{"message": line}
		}
		obj[key] = ts
		b, err := json.Marshal(obj)
		if err != nil {
			return "", fmt.Errorf("transform: marshal: %w", err)
		}
		return string(b), nil
	}
}

// AddField injects a static key/value pair into JSON log lines.
// Non-JSON lines are wrapped in a JSON object first.
func AddField(key, value string) Transformer {
	return func(line string) (string, error) {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			obj = map[string]interface{}{"message": line}
		}
		obj[key] = value
		b, err := json.Marshal(obj)
		if err != nil {
			return "", fmt.Errorf("transform: marshal: %w", err)
		}
		return string(b), nil
	}
}

// Uppercase returns a transformer that upper-cases the entire line.
// Useful mainly for testing and simple non-JSON pipelines.
func Uppercase() Transformer {
	return func(line string) (string, error) {
		return strings.ToUpper(line), nil
	}
}

// PassThrough is a no-op transformer.
func PassThrough() Transformer {
	return func(line string) (string, error) {
		return line, nil
	}
}
