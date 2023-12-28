package health

import (
	"net/http"
	"time"
)

type HealthAggregationsResponse struct {
	InstanceId           string                 `json:"instance_id"`
	IntervalDuration     time.Duration          `json:"interval_duration"`
	IntervalAggregations []*IntervalAggregation `json:"aggregations"`
}

func (s *JsonPollingSink) StartServer(addr string) {
	go http.ListenAndServe(addr, s)
}
