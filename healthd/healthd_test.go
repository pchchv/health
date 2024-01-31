package healthd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pchchv/health"
	"github.com/stretchr/testify/assert"
)

// assertFooBarAggregation asserts that intAgg is the aggregation (generally) of the stuff created in TestHealthD
func assertFooBarAggregation(t *testing.T, intAgg *health.IntervalAggregation) {
	assert.EqualValues(t, intAgg.Events["bar"], 2)
	assert.EqualValues(t, intAgg.Timers["baz"].Count, 2)

	job := intAgg.Jobs["foo"]
	assert.NotNil(t, job)
	assert.EqualValues(t, job.Count, 2)
	assert.EqualValues(t, job.CountSuccess, 1)
	assert.EqualValues(t, job.CountValidationError, 1)
}

func testAggregations(t *testing.T, hd *HealthD) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/healthd/aggregations", nil)
	hd.apiRouter().ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	var resp ApiResponseAggregations
	err := json.Unmarshal(recorder.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, len(resp.Aggregations), 1)
	assertFooBarAggregation(t, resp.Aggregations[0])
}

func testAggregationsOverall(t *testing.T, hd *HealthD) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/healthd/aggregations/overall", nil)
	hd.apiRouter().ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	var resp ApiResponseAggregationsOverall
	err := json.Unmarshal(recorder.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.NotNil(t, resp.Overall)
	assertFooBarAggregation(t, resp.Overall)
}

func testJobs(t *testing.T, hd *HealthD) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/healthd/jobs", nil)
	hd.apiRouter().ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	var resp ApiResponseJobs
	err := json.Unmarshal(recorder.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, len(resp.Jobs), 1)
	job := resp.Jobs[0]
	assert.Equal(t, job.Name, "foo")
	assert.EqualValues(t, job.Count, 2)
	assert.EqualValues(t, job.CountSuccess, 1)
	assert.EqualValues(t, job.CountValidationError, 1)
	assert.EqualValues(t, job.CountError, 0)
	assert.EqualValues(t, job.CountPanic, 0)
	assert.EqualValues(t, job.CountJunk, 0)
	assert.EqualValues(t, job.NanosSum, 14443)
	assert.EqualValues(t, job.NanosMin, 5678)
	assert.EqualValues(t, job.NanosMax, 8765)
	assert.InDelta(t, job.NanosAvg, 7221.5, 0.01)
	assert.InDelta(t, job.NanosSumSquares, 1.09064909e+08, 0.01)
	assert.InDelta(t, job.NanosStdDev, 2182.8386, 0.01)
}
