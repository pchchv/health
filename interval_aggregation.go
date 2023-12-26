package health

type TimerAggregation struct {
	Count           int64   `json:"count"`
	NanosSum        int64   `json:"nanos_sum"`
	NanosSumSquares float64 `json:"nanos_sum_squares"` // 3seconds^2 overflows an int64
	NanosMin        int64   `json:"nanos_min"`
	NanosMax        int64   `json:"nanos_max"`
}
