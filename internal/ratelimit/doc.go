// Package ratelimit implements a token-bucket rate limiter used to cap
// the throughput of log lines flowing through the logpipe pipeline.
//
// Usage:
//
//	limiter := ratelimit.New(ratelimit.Config{
//		Rate:  500,  // 500 lines/sec steady-state
//		Burst: 100,  // allow short bursts up to 100 lines
//	})
//
//	for line := range lines {
//		if err := limiter.Wait(ctx); err != nil {
//			return err
//		}
//		process(line)
//	}
package ratelimit
