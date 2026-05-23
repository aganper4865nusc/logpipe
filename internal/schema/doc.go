// Package schema validates structured log lines against a configurable set of
// required JSON fields.
//
// A Validator is constructed with a list of field names that every log line
// must contain. Lines that are not valid JSON, or that omit one or more
// required fields, are rejected with a descriptive error.
//
// When no required fields are configured the validator acts as a pass-through
// and accepts every line without parsing, keeping the hot path allocation-free.
//
// Typical usage in a pipeline stage:
//
//	v := schema.New([]string{"level", "msg", "ts"})
//	if err := v.Validate(line); err != nil {
//		// drop or flag the line
//	}
package schema
