package health

import "time"

type HealthAggregationsResponse struct {
	InstanceId           string                 `json:"instance_id"`
	IntervalDuration     time.Duration          `json:"interval_duration"`
	IntervalAggregations []*IntervalAggregation `json:"aggregations"`
}
