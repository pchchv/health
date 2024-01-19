package healthd

import (
	"sort"
	"sync"
	"time"

	"github.com/pchchv/health"
)

type HostStatus struct {
	HostPort                string        `json:"host_port"`
	LastCheckTime           time.Time     `json:"last_check_time"`
	LastInstanceId          string        `json:"last_instance_id"`
	LastIntervalDuration    time.Duration `json:"last_interval_duration"`
	LastErr                 string        `json:"last_err"`
	LastNanos               int64         `json:"last_nanos"`
	LastCode                int           `json:"last_code"` // http status code of last response
	FirstSuccessfulResponse time.Time     `json:"first_successful_response"`
	LastSuccessfulResponse  time.Time     `json:"last_successful_response"`
}

type hostAggregationKey struct {
	Time       time.Time
	InstanceId string
	HostPort   string
}

type ByInterval []*health.IntervalAggregation

func (a ByInterval) Len() int {
	return len(a)
}

func (a ByInterval) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByInterval) Less(i, j int) bool {
	return a[i].IntervalStart.Before(a[j].IntervalStart)
}

type HealthD struct {
	stream *health.Stream
	// How long is each aggregation interval, e.g. 1 minute
	intervalDuration time.Duration
	// Retain controls how many metrics interval we keep, e.g. 5 minutes
	retain time.Duration
	// maxIntervals is the maximum length of intervals.
	// It is retain / interval.
	maxIntervals int
	// These guys are the real aggregated deal
	intervalAggregations []*health.IntervalAggregation
	// let's keep the last 5 minutes worth of data from each host
	hostAggregations map[hostAggregationKey]*health.IntervalAggregation
	// intervalsNeedingRecalculation is a set of intervals that need to be recalculated.
	// It is cleared when they are recalculated.
	intervalsNeedingRecalculation map[time.Time]struct{}
	// map from HostPort to status
	hostStatus         map[string]*HostStatus
	intervalsChanChan  chan chan []*health.IntervalAggregation
	hostsChanChan      chan chan []*HostStatus
	stopFlag           int64
	stopAggregator     chan bool
	stopStopAggregator chan bool
	stopHTTP           func() bool
}

// poll is meant to be alled in a new goroutine.
// It will poll each managed host in a new goroutine.
// When everything has finished, it will send nil to responses to signal that we have all data.
func (hd *HealthD) poll(responses chan *pollResponse) {
	var wg sync.WaitGroup
	for _, hs := range hd.hostStatus {
		wg.Add(1)
		go func(hs *HostStatus) {
			defer wg.Done()
			poll(hd.stream, hs.HostPort, responses)
		}(hs)
	}
	wg.Wait()
	responses <- nil
}

// purge purges old hostAggregations older than 5 intervals.
func (agg *HealthD) purge() {
	var threshold = agg.intervalDuration * 5
	for k := range agg.hostAggregations {
		if time.Since(k.Time) > threshold {
			delete(agg.hostAggregations, k)
		}
	}

	n := len(agg.intervalAggregations)
	if n > agg.maxIntervals {
		agg.intervalAggregations = agg.intervalAggregations[(n - agg.maxIntervals):]
	}
}

func (agg *HealthD) setAggregation(intAgg *health.IntervalAggregation) {
	// If we already have the intAgg, replace it.
	for i, existingAgg := range agg.intervalAggregations {
		if existingAgg.IntervalStart == intAgg.IntervalStart {
			agg.intervalAggregations[i] = intAgg
			return
		}
	}

	// Otherwise, just append it and sort to get ordering right.
	agg.intervalAggregations = append(agg.intervalAggregations, intAgg)
	sort.Sort(ByInterval(agg.intervalAggregations))

	// If we have too many aggregations, truncate
	n := len(agg.intervalAggregations)
	if n > agg.maxIntervals {
		agg.intervalAggregations = agg.intervalAggregations[(n - agg.maxIntervals):]
	}
}

func (hd *HealthD) recalculateIntervals() {
	job := hd.stream.NewJob("recalculate")
	for k := range hd.intervalsNeedingRecalculation {
		intAggsAtTime := []*health.IntervalAggregation{}
		for key, intAgg := range hd.hostAggregations {
			if key.Time == k {
				intAggsAtTime = append(intAggsAtTime, intAgg)
			}
		}

		overallAgg := health.NewIntervalAggregation(k)
		for _, ia := range intAggsAtTime {
			overallAgg.Merge(ia)
		}
		hd.setAggregation(overallAgg)
	}

	// reset everything:
	hd.intervalsNeedingRecalculation = make(map[time.Time]struct{})
	job.Complete(health.Success)
}
