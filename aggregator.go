package health

import "time"

var nowMock time.Time

type aggregator struct {
	// How long is each aggregation interval. Eg, 1 minute
	intervalDuration time.Duration
	// Retain controls how many metrics interval we keep. Eg, 5 minutes
	retain time.Duration
	// maxIntervals is the maximum length of intervals.
	// It is retain / interval.
	maxIntervals int
	// intervals is a slice of the retained intervals
	intervalAggregations []*IntervalAggregation
}

func newAggregator(intervalDuration time.Duration, retain time.Duration) *aggregator {
	maxIntervals := int(retain / intervalDuration)
	return &aggregator{
		intervalDuration:     intervalDuration,
		retain:               retain,
		maxIntervals:         maxIntervals,
		intervalAggregations: make([]*IntervalAggregation, 0, maxIntervals),
	}
}

func now() time.Time {
	if nowMock.IsZero() {
		return time.Now()
	}
	return nowMock
}
