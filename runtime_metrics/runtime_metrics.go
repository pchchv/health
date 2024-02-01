package runtime_metrics

import "time"

type Options struct {
	Interval   time.Duration
	Memory     bool
	GC         bool
	GCQuantile bool
	Goroutines bool
	Cgo        bool
	FDs        bool
}
