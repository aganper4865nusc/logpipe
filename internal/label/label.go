// Package label provides a transform that attaches static key-value labels
// to every log line that passes through the pipeline.
package label

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Labeler attaches a fixed set of labels to each log line.
type Labeler struct {
	labels map[string]string
}

// New creates a Labeler that will inject the given labels into every line.
// Labels must be non-empty; an error is returned if the map is nil or empty.
func New(labels map[string]string) (*Labeler, error) {
	if len(labels) == 0 {
		return nil, fmt.Errorf("label: at least one label is required")
	}
	copy := make(map[string]string, len(labels))
	for k, v := range labels {
		if strings.TrimSpace(k) == "" {
			return nil, fmt.Errorf("label: key must not be blank")
		}
		copy[k] = v
	}
	return &Labeler{labels: copy}, nil
}

// Transform injects the configured labels into line.
// If line is valid JSON, labels are merged into the object.
// Otherwise the labels are appended as key=value pairs separated by spaces.
func (l *Labeler) Transform(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err == nil {
		for k, v := range l.labels {
			if _, exists := obj[k]; !exists {
				obj[k] = v
			}
		}
		b, err := json.Marshal(obj)
		if err == nil {
			return string(b)
		}
	}
	var sb strings.Builder
	sb.WriteString(line)
	for k, v := range l.labels {
		sb.WriteByte(' ')
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(v)
	}
	return sb.String()
}

// PassThrough returns a transform function that applies no labels (identity).
func PassThrough(line string) string { return line }
