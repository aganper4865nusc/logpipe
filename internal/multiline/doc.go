// Package multiline implements a multi-line log assembler for logpipe.
//
// Many log formats (Java stack traces, Python tracebacks, multi-line JSON)
// span several raw lines that logically belong to a single event. This
// package detects event boundaries using a configurable start-pattern
// regular expression and folds continuation lines into one string before
// the event is forwarded to the pipeline.
//
// Basic usage:
//
//	a, err := multiline.New(multiline.Config{
//		StartPattern: `^\d{4}-\d{2}-\d{2}`,  // ISO-date starts a new event
//		MaxLines:     100,
//		FlushTimeout: 2 * time.Second,
//	})
//	if err != nil { /* handle */ }
//
//	for _, raw := range rawLines {
//		if event, ready := a.Feed(raw); ready {
//			forward(event)
//		}
//	}
//	if tail := a.Flush(); tail != "" {
//		forward(tail)
//	}
package multiline
