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
