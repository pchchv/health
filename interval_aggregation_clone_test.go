package health

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type eventErr struct {
	event string
	err   error
}

func assertAggregationData(t *testing.T, intAgg *IntervalAggregation) {
	assert.Equal(t, 300, len(intAgg.Jobs))
	assert.Equal(t, 1200, len(intAgg.Events))
	assert.Equal(t, 1200, len(intAgg.Timers))
	assert.Equal(t, 1200, len(intAgg.Gauges))
	assert.Equal(t, 1200, len(intAgg.EventErrs))

	// Spot-check events:
	assert.EqualValues(t, 1, intAgg.Events["event0"])

	// Spot check gauges:
	assert.EqualValues(t, 3.14, intAgg.Gauges["gauge0"])

	// Spot-check timings:
	assert.EqualValues(t, 1, intAgg.Timers["timing0"].Count)
	assert.EqualValues(t, 12, intAgg.Timers["timing0"].NanosSum)

	// Spot-check event-errs:
	assert.EqualValues(t, 1, intAgg.EventErrs["err0"].Count)
	assert.Equal(t, []error{errors.New("wat")}, intAgg.EventErrs["err0"].getErrorSamples())

	// Spot-check jobs:
	job := intAgg.Jobs["job0"]
	assert.EqualValues(t, 1, job.CountSuccess)
	assert.EqualValues(t, 0, job.CountError)
	assert.EqualValues(t, 1, job.Events["event0"])
	assert.EqualValues(t, 0, job.Events["event4"])
	assert.EqualValues(t, 3.14, job.Gauges["gauge0"])
	assert.EqualValues(t, 0.0, job.Gauges["gauge4"])
	assert.EqualValues(t, 1, job.Timers["timing0"].Count)
	assert.EqualValues(t, 12, job.Timers["timing0"].NanosSum)
	assert.EqualValues(t, 1, job.EventErrs["err0"].Count)
	assert.Equal(t, []error{errors.New("wat")}, job.EventErrs["err0"].getErrorSamples())

	// Nothing foo or bar related
	_, ok := intAgg.Jobs["foo"]
	assert.False(t, ok)
	assert.EqualValues(t, 0, intAgg.Events["bar"])
	assert.Nil(t, intAgg.Timers["bar"])
	assert.Nil(t, intAgg.EventErrs["bar"])
}
