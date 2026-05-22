package pipeline

import (
	"log"

	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/sink"
)

// Sink is the interface satisfied by all output sinks.
type Sink interface {
	Write(line string) error
}

// Pipeline connects a line source to one or more sinks, applying a filter.
type Pipeline struct {
	filter *filter.Filter
	sinks  []Sink
	lines  <-chan string
}

// Config holds construction parameters for a Pipeline.
type Config struct {
	Filter *filter.Filter
	Sinks  []Sink
	Lines  <-chan string
}

// New creates a Pipeline from the provided Config.
func New(cfg Config) *Pipeline {
	f := cfg.Filter
	if f == nil {
		f = filter.PassThrough()
	}
	return &Pipeline{
		filter: f,
		sinks:  cfg.Sinks,
		lines:  cfg.Lines,
	}
}

// Run reads from the line channel, filters, and forwards to all sinks.
// It blocks until the channel is closed.
func (p *Pipeline) Run() {
	for line := range p.lines {
		if !p.filter.Match(line) {
			continue
		}
		for _, s := range p.sinks {
			if err := s.Write(line); err != nil {
				log.Printf("pipeline: sink write error: %v", err)
			}
		}
	}
}

// Ensure *sink.StdoutSink and *sink.FileSink satisfy Sink at compile time.
var _ Sink = (*sink.StdoutSink)(nil)
var _ Sink = (*sink.FileSink)(nil)
