// Package sampling implements probabilistic log-line sampling for logpipe.
//
// A [Sampler] wraps a source pipeline stage and stochastically drops lines
// according to a configured rate in the range (0, 1].  A rate of 1.0 is a
// no-op (all lines pass through), while 0.1 forwards approximately 10 % of
// observed lines.
//
// Usage:
//
//	s, err := sampling.New(sampling.Config{Rate: 0.25, Seed: time.Now().UnixNano()})
//	if err != nil { ... }
//	if s.Sample(line) {
//	    // forward line to sink
//	}
//
// For pipelines that do not require sampling, use [PassThrough] which always
// returns true and incurs negligible overhead.
package sampling
