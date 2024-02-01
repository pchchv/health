package runtime_metrics

import (
	"time"

	"github.com/pchchv/health"
)

type Options struct {
	Interval   time.Duration
	Memory     bool
	GC         bool
	GCQuantile bool
	Goroutines bool
	Cgo        bool
	FDs        bool
}

type RuntimeMetrics struct {
	stream       health.EventReceiver
	options      Options
	stopChan     chan bool
	stopStopChan chan bool
}
