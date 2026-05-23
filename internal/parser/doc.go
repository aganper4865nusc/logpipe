// Package parser provides composable log line parsers for logpipe.
//
// Three built-in formats are supported:
//
//   - JSON: unmarshals a JSON object from each line.
//   - Regex: extracts named capture groups from each line.
//   - KV: splits lines into key=value pairs using configurable separators.
//
// A PassThrough parser is also provided for pipelines that do not require
// structured parsing; it wraps the raw line under the "message" key.
//
// All parsers implement the ParseFn signature:
//
//	type ParseFn func(line string) (map[string]any, error)
//
// Errors are returned when a line cannot be parsed according to the chosen
// format; callers should decide whether to drop or forward the raw line.
package parser
