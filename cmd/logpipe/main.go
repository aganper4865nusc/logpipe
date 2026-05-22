// main is the entry point for the logpipe log aggregator.
// It loads configuration, initializes tailers and sinks, and
// routes log lines from sources to their configured destinations.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/yourorg/logpipe/internal/config"
	"github.com/yourorg/logpipe/internal/sink"
	"github.com/yourorg/logpipe/internal/tail"
)

func main() {
	cfgPath := flag.String("config", "logpipe.yaml", "path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Build all configured sinks.
	sinks := make([]sink.Sink, 0, len(cfg.Sinks))
	for _, sc := range cfg.Sinks {
		s, err := sink.New(sc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create sink %q: %v\n", sc.Type, err)
			os.Exit(1)
		}
		sinks = append(sinks, s)
	}

	// lines is a shared channel that all tailers write into.
	lines := make(chan string, cfg.Workers*10)

	var tailerWg sync.WaitGroup

	// Start a tailer goroutine for each configured source.
	for _, src := range cfg.Sources {
		t, err := tail.New(src.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to tail %q: %v\n", src.Path, err)
			os.Exit(1)
		}

		tailerWg.Add(1)
		go func(t *tail.Tailer) {
			defer tailerWg.Done()
			for line := range t.Lines() {
				lines <- line
			}
		}(t)
	}

	// Close the lines channel once all tailers have finished.
	go func() {
		tailerWg.Wait()
		close(lines)
	}()

	// Start worker goroutines that fan out each line to every sink.
	var workerWg sync.WaitGroup
	for i := 0; i < cfg.Workers; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for line := range lines {
				for _, s := range sinks {
					if err := s.Write(line); err != nil {
						log.Printf("sink write error: %v", err)
					}
				}
			}
		}()
	}

	// Block until SIGINT or SIGTERM is received.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down logpipe…")

	// Workers will drain and exit once the lines channel is closed.
	workerWg.Wait()

	// Close all sinks to flush any buffered data.
	for _, s := range sinks {
		if err := s.Close(); err != nil {
			log.Printf("sink close error: %v", err)
		}
	}
}
