// Package dedupe implements a sliding-window deduplication layer for log lines.
//
// A Deduplicator records each unique line it sees and suppresses subsequent
// occurrences of the same line that arrive within a configurable time window.
// Once the window elapses the line is forgotten and will be forwarded again if
// it reappears.
//
// Example usage:
//
//	dd := dedupe.New(10 * time.Second)
//	for line := range lines {
//		if !dd.IsDuplicate(line) {
//			sink.Write(line)
//		}
//	}
//
// Setting the window to zero disables deduplication entirely, making every
// call to IsDuplicate return false.
package dedupe
