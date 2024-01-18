package healthd

import "time"

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
