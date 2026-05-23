// Package normalize provides composable TransformFunc helpers that clean and
// standardise raw log lines before they enter the processing pipeline.
//
// Available transforms:
//
//   - PassThrough  – identity, returns the line unchanged.
//   - TrimSpace    – strips leading/trailing whitespace.
//   - CollapseWhitespace – collapses internal whitespace runs to a single space.
//   - StripControl – removes non-printable control characters (tabs are kept).
//   - NormalizeNewlines – converts CR+LF and bare CR to LF.
//
// Transforms can be composed with Chain:
//
//	fn := normalize.Chain(
//	    normalize.NormalizeNewlines,
//	    normalize.TrimSpace,
//	    normalize.StripControl,
//	)
//	cleaned := fn(rawLine)
package normalize
