// Package throttle wraps a Sink function and enforces a maximum throughput
// expressed as lines per second.
//
// A token-bucket algorithm is used: a background goroutine refills the bucket
// at the configured rate. Write blocks until a token is available or the
// supplied context is cancelled.
//
// Example usage:
//
//	th, cancel, err := throttle.New(mySink, 200) // 200 lines/s
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer cancel()
//
//	if err := th.Write(ctx, line); err != nil {
//		log.Println("write failed:", err)
//	}
package throttle
