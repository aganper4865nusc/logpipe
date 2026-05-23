// Package normalize provides line normalization transforms that
// standardize whitespace, encoding, and line endings before further
// processing by the pipeline.
package normalize

import (
	"strings"
	"unicode"
)

// TransformFunc transforms a single log line and returns the result.
type TransformFunc func(line string) string

// Chain returns a TransformFunc that applies each fn in order.
func Chain(fns ...TransformFunc) TransformFunc {
	return func(line string) string {
		for _, fn := range fns {
			line = fn(line)
		}
		return line
	}
}

// PassThrough returns the line unchanged.
func PassThrough(line string) string { return line }

// TrimSpace removes leading and trailing whitespace (spaces, tabs, newlines).
func TrimSpace(line string) string {
	return strings.TrimSpace(line)
}

// CollapseWhitespace replaces runs of internal whitespace with a single space.
func CollapseWhitespace(line string) string {
	var b strings.Builder
	b.Grow(len(line))
	inSpace := false
	for _, r := range line {
		if unicode.IsSpace(r) {
			if !inSpace {
				b.WriteRune(' ')
				inSpace = true
			}
		} else {
			b.WriteRune(r)
			inSpace = false
		}
	}
	return b.String()
}

// StripControl removes non-printable control characters (except tab).
func StripControl(line string) string {
	var b strings.Builder
	b.Grow(len(line))
	for _, r := range line {
		if r == '\t' || unicode.IsPrint(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// NormalizeNewlines replaces CR+LF and bare CR with LF.
func NormalizeNewlines(line string) string {
	line = strings.ReplaceAll(line, "\r\n", "\n")
	line = strings.ReplaceAll(line, "\r", "\n")
	return line
}
