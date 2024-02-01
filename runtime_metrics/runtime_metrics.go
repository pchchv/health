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

func NewRuntimeMetrics(stream health.EventReceiver, options *Options) *RuntimeMetrics {
	rm := &RuntimeMetrics{
		stream:       stream,
		stopChan:     make(chan bool),
		stopStopChan: make(chan bool),
	}
	if options != nil {
		rm.options = *options
	} else {
		rm.options = Options{time.Second * 5, true, true, true, true, true, true}
	}
	return rm
}
