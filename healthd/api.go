package healthd

import (
	"encoding/json"
	"fmt"
	"net/http"
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
