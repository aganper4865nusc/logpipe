// Package masker provides field-level value masking for structured (JSON)
// log lines.
//
// Usage:
//
//	masker := masker.New([]string{"password", "token"}, "***")
//	sanitised := masker.Transform(rawLine)
//
// Non-JSON lines are passed through without modification. The mask string
// defaults to "***" when an empty string is supplied to New.
//
// PassThrough is provided as a convenience no-op transform for use when
// masking is disabled by configuration.
package masker
