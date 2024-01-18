package healthd

import (
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
