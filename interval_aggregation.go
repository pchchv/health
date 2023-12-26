package health

import "reflect"

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
	// Let's keep a ring buffer of some errors. I feel like this isn't the best data structure / plan of attack here but works for now.
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
	} else if !reflect.DeepEqual(lastErr, inputErr) {
		n := len(ec.errorSamples)
		ec.errorSampleIndex = (ec.errorSampleIndex + 1) % n
		ec.errorSamples[ec.errorSampleIndex] = inputErr
	}
}
