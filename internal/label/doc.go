// Package label provides a Transform function that attaches static key-value
// labels to every log line flowing through a logpipe pipeline.
//
// For JSON lines the labels are merged into the top-level object without
// overwriting keys that are already present.  For plain-text lines the labels
// are appended as space-separated key=value pairs.
//
// Example usage:
//
//	l, err := label.New(map[string]string{
//		"env":     "production",
//		"service": "api",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	output := l.Transform(inputLine)
package label
