package healthd

import (
	"time"

	"github.com/pchchv/health"
)

type pollResponse struct {
	HostPort  string
	Timestamp time.Time
	Err       error
	Code      int
	Nanos     int64
	health.HealthAggregationsResponse
}

func metricsUrl(hostPort string) string {
	return "http://" + hostPort + "/health"
}
