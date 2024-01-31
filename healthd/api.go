package healthd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pchchv/grom"
	"github.com/pchchv/health"
)

type apiResponse struct {
	InstanceId       string        `json:"instance_id"`
	IntervalDuration time.Duration `json:"interval_duration"`
}

type ApiResponseAggregations struct {
	apiResponse
	Aggregations []*health.IntervalAggregation `json:"aggregations"`
}

type ApiResponseAggregationsOverall struct {
	apiResponse
	Overall *health.IntervalAggregation `json:"overall"`
}

type ApiResponseHosts struct {
	apiResponse
	Hosts []*HostStatus `json:"hosts"`
}

// Job represents a health.JobAggregation,
// but designed for JSON-ization without all the nested counters/timers
type Job struct {
	Name                 string  `json:"name"`
	Count                int64   `json:"count"`
	CountSuccess         int64   `json:"count_success"`
	CountValidationError int64   `json:"count_validation_error"`
	CountPanic           int64   `json:"count_panic"`
	CountError           int64   `json:"count_error"`
	CountJunk            int64   `json:"count_junk"`
	NanosSum             int64   `json:"nanos_sum"`
	NanosSumSquares      float64 `json:"nanos_sum_squares"`
	NanosMin             int64   `json:"nanos_min"`
	NanosMax             int64   `json:"nanos_max"`
	NanosAvg             float64 `json:"nanos_avg"`
	NanosStdDev          float64 `json:"nanos_std_dev"`
}

type ApiResponseJobs struct {
	apiResponse
	Jobs []*Job `json:"jobs"`
}

type apiContext struct {
	hd *HealthD
	*health.Job
}

func (c *apiContext) SetContentType(rw grom.ResponseWriter, req *grom.Request, next grom.NextMiddlewareFunc) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	next(rw, req)
}

func (c *apiContext) HealthMiddleware(rw grom.ResponseWriter, r *grom.Request, next grom.NextMiddlewareFunc) {
	c.Job = c.hd.stream.NewJob(r.RoutePath())

	path := r.URL.Path
	c.EventKv("starting_request", health.Kvs{"path": path})

	next(rw, r)

	code := rw.StatusCode()
	kvs := health.Kvs{
		"code": fmt.Sprint(code),
		"path": path,
	}

	// Map HTTP status code to category.
	var status health.CompletionStatus
	if code < 400 {
		status = health.Success
	} else if code < 500 {
		if code == 422 {
			status = health.ValidationError
		} else {
			status = health.Junk // 404, 401
		}
	} else {
		status = health.Error
	}
	c.CompleteKv(status, kvs)
}

func (c *apiContext) Aggregations(rw grom.ResponseWriter, r *grom.Request) {
	aggregations := c.hd.getAggregationSequence()
	resp := &ApiResponseAggregations{
		apiResponse:  getApiResponse(c.hd.intervalDuration),
		Aggregations: aggregations,
	}
	renderJson(rw, resp)
}

func (c *apiContext) Overall(rw grom.ResponseWriter, r *grom.Request) {
	aggregations := c.hd.getAggregationSequence()
	overall := combineAggregations(aggregations)
	resp := &ApiResponseAggregationsOverall{
		apiResponse: getApiResponse(c.hd.intervalDuration),
		Overall:     overall,
	}
	renderJson(rw, resp)
}

// By is the type of a "less" function
// that defines the ordering of its Planet arguments.
type By func(j1, j2 *Job) bool

type jobSorter struct {
	jobs []*Job
	by   By
}

func getApiResponse(duration time.Duration) apiResponse {
	return apiResponse{
		InstanceId:       health.Identifier,
		IntervalDuration: duration,
	}
}

func renderNotFound(rw http.ResponseWriter) {
	rw.WriteHeader(404)
	fmt.Fprintf(rw, `{"error": "not_found"}`)
}

func renderError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(500)
	fmt.Fprintf(rw, `{"error": "%s"}`, err.Error())
}

func renderJson(rw http.ResponseWriter, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		renderError(rw, err)
		return
	}
	fmt.Fprintf(rw, string(jsonData))
}

func combineAggregations(aggregations []*health.IntervalAggregation) *health.IntervalAggregation {
	if len(aggregations) == 0 {
		return nil
	}

	overallAgg := health.NewIntervalAggregation(aggregations[0].IntervalStart)
	for _, ia := range aggregations {
		overallAgg.Merge(ia)
	}
	return overallAgg
}

func getSort(r *grom.Request) string {
	return r.URL.Query().Get("sort")
}

func getLimit(r *grom.Request) int {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		return 0
	}

	n, err := strconv.ParseInt(limit, 10, 0)
	if err != nil {
		return 0
	}
	return int(n)
}
