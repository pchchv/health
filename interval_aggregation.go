package health

import (
	"errors"
	"time"
)

type TimerAggregation struct {
	Count           int64   `json:"count"`
	NanosSum        int64   `json:"nanos_sum"`
	NanosSumSquares float64 `json:"nanos_sum_squares"` // 3seconds^2 overflows an int64
	NanosMin        int64   `json:"nanos_min"`
	NanosMax        int64   `json:"nanos_max"`
}

func (a *TimerAggregation) ingest(nanos int64) {
	a.Count++
	a.NanosSum += nanos
	floatNano := float64(nanos)
	a.NanosSumSquares += (floatNano * floatNano)
	if a.Count == 1 || nanos < a.NanosMin {
		a.NanosMin = nanos
	}

	if a.Count == 1 || nanos > a.NanosMax {
		a.NanosMax = nanos
	}
}

type ErrorCounter struct {
	Count int64 `json:"count"`
	// Let's keep a ring buffer of some errors.
	// I feel like this isn't the best data structure / plan of attack here but works for now.
	errorSamples     [5]error
	errorSampleIndex int
}

func (ec *ErrorCounter) incrementAndAddError(inputErr error) {
	ec.Count++
	ec.addError(inputErr)
}

func (ec *ErrorCounter) addError(inputErr error) {
	lastErr := ec.errorSamples[ec.errorSampleIndex]
	if lastErr == nil {
		ec.errorSamples[ec.errorSampleIndex] = inputErr
	} else if !errors.Is(lastErr, inputErr) {
		n := len(ec.errorSamples)
		ec.errorSampleIndex = (ec.errorSampleIndex + 1) % n
		ec.errorSamples[ec.errorSampleIndex] = inputErr
	}
}

func (ec *ErrorCounter) getErrorSamples() []error {
	// count how many non-nil errors are there so we can make a slice of the right size
	count := 0
	for _, e := range ec.errorSamples {
		if e != nil {
			count++
		}
	}
	ret := make([]error, 0, count)

	// put non-nil errors in slice
	for _, e := range ec.errorSamples {
		if e != nil {
			ret = append(ret, e)
		}
	}
	return ret
}

type aggregationMaps struct {
	Timers    map[string]*TimerAggregation `json:"timers"`
	Gauges    map[string]float64           `json:"gauges"`
	Events    map[string]int64             `json:"events"`
	EventErrs map[string]*ErrorCounter     `json:"event_errs"`
}

func (am *aggregationMaps) initAggregationMaps() {
	am.Timers = make(map[string]*TimerAggregation)
	am.Gauges = make(map[string]float64)
	am.Events = make(map[string]int64)
	am.EventErrs = make(map[string]*ErrorCounter)
}

func (am *aggregationMaps) getCounterErrs(event string) *ErrorCounter {
	ce := am.EventErrs[event]
	if ce == nil {
		ce = &ErrorCounter{}
		am.EventErrs[event] = ce
	}
	return ce
}

func (am *aggregationMaps) getTimers(event string) *TimerAggregation {
	t := am.Timers[event]
	if t == nil {
		t = &TimerAggregation{}
		am.Timers[event] = t
	}
	return t
}

type JobAggregation struct {
	aggregationMaps
	TimerAggregation
	CountValidationError int64 `json:"count_validation_error"`
	CountSuccess         int64 `json:"count_success"`
	CountPanic           int64 `json:"count_panic"`
	CountError           int64 `json:"count_error"`
	CountJunk            int64 `json:"count_junk"`
}

func (a *JobAggregation) ingest(status CompletionStatus, nanos int64) {
	a.TimerAggregation.ingest(nanos)
	switch status {
	case Success:
		a.CountSuccess++
	case ValidationError:
		a.CountValidationError++
	case Panic:
		a.CountPanic++
	case Error:
		a.CountError++
	case Junk:
		a.CountJunk++
	}
}

// IntervalAggregation will hold data for a given aggregation interval.
type IntervalAggregation struct {
	// aggregationMaps will hold event/timer information that is not nested per-job.
	aggregationMaps
	// The start time of the interval
	IntervalStart time.Time `json:"interval_start"`
	// SerialNumber increments every time the aggregation changes.
	// It does not increment if the aggregation does not change.
	SerialNumber int64 `json:"serial_number"`
	// Jobs hold a map of job name -> data about that job.
	// This includes both primary-job information (success vs error, et all) as well as
	// scoping timers/counters by the job.
	Jobs map[string]*JobAggregation `json:"jobs"`
}

func NewIntervalAggregation(intervalStart time.Time) *IntervalAggregation {
	intAgg := &IntervalAggregation{
		IntervalStart: intervalStart,
		Jobs:          make(map[string]*JobAggregation),
	}
	intAgg.initAggregationMaps()
	return intAgg
}

func (ia *IntervalAggregation) getJobAggregation(job string) *JobAggregation {
	jobAgg := ia.Jobs[job]
	if jobAgg == nil {
		jobAgg = &JobAggregation{}
		jobAgg.initAggregationMaps()
		ia.Jobs[job] = jobAgg
	}
	return jobAgg
}
