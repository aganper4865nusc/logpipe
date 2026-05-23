// Package fanout provides a concurrent fan-out writer that delivers each log
// line to an arbitrary number of Sink implementations simultaneously.
//
// Usage:
//
//	f := fanout.New(sinkA, sinkB)
//	f.Add(sinkC)          // register additional sinks at any time
//	err := f.Write(line)  // delivered to all three sinks concurrently
//
// If one or more sinks return an error, Write returns a *MultiError containing
// every individual failure; successful sinks are unaffected.
package fanout
