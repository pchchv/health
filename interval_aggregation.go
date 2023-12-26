package health

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
