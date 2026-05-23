// Package checkpoint provides durable offset tracking for log sources.
//
// When logpipe restarts it can resume tailing each file from the last
// committed byte offset rather than re-processing lines from the
// beginning. Offsets are written atomically via a rename so a crash
// during a flush never leaves the checkpoint file in a corrupt state.
//
// Usage:
//
//	store, err := checkpoint.New("/var/lib/logpipe/offsets.json")
//	if err != nil { ... }
//
//	// Restore the previous position before starting the tailer.
//	offset := store.Get("app.log")
//
//	// After processing each line, persist the new position.
//	_ = store.Set("app.log", newOffset)
package checkpoint
