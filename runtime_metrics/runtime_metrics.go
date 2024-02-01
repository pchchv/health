package runtime_metrics

import (
	"runtime"
	"runtime/debug"
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

func (rm *RuntimeMetrics) Report() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	if rm.options.Memory {
		// bytes allocated and not yet freed
		rm.reportGauge("alloc", float64(mem.Alloc))
		// total number of allocated objects
		rm.reportGauge("heap_objects", float64(mem.HeapObjects))
	}

	if rm.options.GC {
		rm.reportGauge("pause_total_ns", float64(mem.PauseTotalNs))
		rm.reportGauge("num_gc", float64(mem.NumGC))
		rm.reportGauge("next_gc", float64(mem.NextGC))
		rm.reportGauge("gc_cpu_fraction", mem.GCCPUFraction)
	}

	if rm.options.GCQuantile {
		var gc debug.GCStats
		gc.PauseQuantiles = make([]time.Duration, 3)
		debug.ReadGCStats(&gc)
		rm.reportGauge("gc_pause_quantile_50", float64(gc.PauseQuantiles[1]/1000)/1000.0)
		rm.reportGauge("gc_pause_quantile_max", float64(gc.PauseQuantiles[2]/1000)/1000.0)
	}

	if rm.options.Goroutines {
		rm.reportGauge("num_goroutines", float64(runtime.NumGoroutine()))
	}

	if rm.options.Cgo {
		rm.reportGauge("num_cgo_call", float64(runtime.NumCgoCall()))
	}

	if rm.options.FDs {
		if num, err := getFDUsage(); err == nil {
			rm.reportGauge("num_fds_used", float64(num))
		}
	}
}

func (rm *RuntimeMetrics) reportGauge(key string, val float64) {
	rm.stream.Gauge(key, val)
}
