package health

import (
	"fmt"
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

func renderNotFound(rw http.ResponseWriter) {
	rw.WriteHeader(404)
	fmt.Fprintf(rw, `{"error": "not_found"}`)
}

func renderError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(500)
	fmt.Fprintf(rw, `{"error": "%s"}`, err.Error())
}
