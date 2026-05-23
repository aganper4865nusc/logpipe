// Package routing implements source-to-sink routing for logpipe.
//
// A [Router] is constructed from a list of [Route] values, each binding
// a named log source (e.g. a tailed file or stdin) to one or more named
// sinks (e.g. "stdout", "file-error").
//
// Routes can also be added at runtime via [Router.Add], making it
// straightforward to wire up dynamic pipeline topologies.
//
// Typical usage:
//
//	router, err := routing.New([]routing.Route{
//		{Source: "app-log",   Sinks: []string{"stdout", "s3"}},
//		{Source: "audit-log", Sinks: []string{"file"}},
//	})
//	if err != nil { ... }
//
//	sinks, ok := router.Resolve("app-log")
//	// sinks == ["stdout", "s3"]
package routing
