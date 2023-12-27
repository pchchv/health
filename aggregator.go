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

func (a *aggregator) createIntervalAggregation(interval time.Time) *IntervalAggregation {
	// Make new interval:
	current := NewIntervalAggregation(interval)
	// If we've reached our max intervals,
	// and we're going to shift everything down,
	// then set the last one
	n := len(a.intervalAggregations)
	if n == a.maxIntervals {
		for i := 1; i < n; i++ {
			a.intervalAggregations[i-1] = a.intervalAggregations[i]
		}
		a.intervalAggregations[n-1] = current
	} else {
		a.intervalAggregations = append(a.intervalAggregations, current)
	}
	return current
}

func (a *aggregator) getIntervalAggregation() *IntervalAggregation {
	intervalStart := now().Truncate(a.intervalDuration)
	n := len(a.intervalAggregations)
	if n > 0 && a.intervalAggregations[n-1].IntervalStart == intervalStart {
		return a.intervalAggregations[n-1]
	}
	return a.createIntervalAggregation(intervalStart)
}

func now() time.Time {
	if nowMock.IsZero() {
		return time.Now()
	}
	return nowMock
}
