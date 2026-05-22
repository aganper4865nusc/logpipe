// Package metrics provides lightweight runtime counters for logpipe.
//
// A Registry tracks per-source and global line counts across four dimensions:
//   - LinesRead      — lines ingested from a tailed source
//   - LinesFiltered  — lines dropped by the filter stage
//   - LinesForwarded — lines successfully delivered to a sink
//   - SinkErrors     — delivery failures reported by a sink
//
// Usage:
//
//	reg := metrics.New()
//	reg.RecordRead("app.log")
//	reg.RecordForwarded("app.log")
//	snap := reg.Snapshot()
//
// A Reporter can be started in a goroutine to emit periodic summaries:
//
//	rep := metrics.NewReporter(reg, os.Stderr, 30*time.Second)
//	go rep.Run(ctx)
package metrics
