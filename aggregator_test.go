package health

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAggregator(t *testing.T) {
	a := newAggregator(time.Minute, time.Minute*5)
	assert.Equal(t, time.Minute, a.intervalDuration)
	assert.Equal(t, time.Minute*5, a.retain)
	assert.Equal(t, 5, a.maxIntervals)
	assert.Equal(t, 0, len(a.intervalAggregations))
	assert.NotNil(t, a.intervalAggregations)
}

func TestEmitEvent(t *testing.T) {
	// Set time, and do a single event
	setNowMock("2011-09-09T23:36:13Z")
	defer resetNowMock()
	a := newAggregator(time.Minute, time.Minute*5)
	a.EmitEvent("foo", "bar")

	assert.Equal(t, 1, len(a.intervalAggregations))

	intAgg := a.intervalAggregations[0]
	assert.NotNil(t, intAgg.Events)
	assert.EqualValues(t, 1, intAgg.Events["bar"])
	assert.EqualValues(t, 1, intAgg.SerialNumber)

	assert.NotNil(t, intAgg.Jobs)
	jobAgg := intAgg.Jobs["foo"]
	assert.NotNil(t, jobAgg)
	assert.NotNil(t, jobAgg.Events)
	assert.EqualValues(t, 1, jobAgg.Events["bar"])

	// Now, without changing the time, we'll do 3 more events:
	a.EmitEvent("foo", "bar") // duplicate to above
	a.EmitEvent("foo", "baz") // same job, diff event
	a.EmitEvent("wat", "bar") // diff job, same event

	assert.Equal(t, 1, len(a.intervalAggregations))

	intAgg = a.intervalAggregations[0]
	assert.EqualValues(t, 3, intAgg.Events["bar"])
	assert.EqualValues(t, 4, intAgg.SerialNumber)

	jobAgg = intAgg.Jobs["foo"]
	assert.EqualValues(t, 2, jobAgg.Events["bar"])
	assert.EqualValues(t, 1, jobAgg.Events["baz"])

	jobAgg = intAgg.Jobs["wat"]
	assert.NotNil(t, jobAgg)
	assert.EqualValues(t, 1, jobAgg.Events["bar"])

	// increment time and do one more event:
	setNowMock("2011-09-09T23:37:01Z")
	a.EmitEvent("foo", "bar")

	assert.Equal(t, 2, len(a.intervalAggregations))

	// make sure old values don't change:
	intAgg = a.intervalAggregations[0]
	assert.EqualValues(t, 3, intAgg.Events["bar"])
	assert.EqualValues(t, 4, intAgg.SerialNumber)

	intAgg = a.intervalAggregations[1]
	assert.EqualValues(t, 1, intAgg.Events["bar"])
	assert.EqualValues(t, 1, intAgg.SerialNumber)
}

func TestEmitEventErr(t *testing.T) {
	setNowMock("2011-09-09T23:36:13Z")
	defer resetNowMock()
	a := newAggregator(time.Minute, time.Minute*5)
	a.EmitEventErr("foo", "bar", errors.New("wat"))

	assert.Equal(t, 1, len(a.intervalAggregations))

	intAgg := a.intervalAggregations[0]
	assert.NotNil(t, intAgg.EventErrs)
	ce := intAgg.EventErrs["bar"]
	assert.NotNil(t, ce)
	assert.EqualValues(t, 1, ce.Count)
	assert.Equal(t, []error{errors.New("wat")}, ce.getErrorSamples())
	assert.EqualValues(t, 1, intAgg.SerialNumber)

	assert.NotNil(t, intAgg.Jobs)
	jobAgg := intAgg.Jobs["foo"]
	assert.NotNil(t, jobAgg)
	assert.NotNil(t, jobAgg.EventErrs)
	ce = jobAgg.EventErrs["bar"]
	assert.EqualValues(t, 1, ce.Count)
	assert.Equal(t, []error{errors.New("wat")}, ce.getErrorSamples())

	// One more event with the same error:
	a.EmitEventErr("foo", "bar", errors.New("wat"))

	intAgg = a.intervalAggregations[0]
	ce = intAgg.EventErrs["bar"]
	assert.EqualValues(t, 2, ce.Count)
	assert.Equal(t, []error{errors.New("wat")}, ce.getErrorSamples()) // doesn't change

	// One more event with diff error:
	a.EmitEventErr("foo", "bar", errors.New("lol"))

	intAgg = a.intervalAggregations[0]
	ce = intAgg.EventErrs["bar"]
	assert.EqualValues(t, 3, ce.Count)
	assert.Equal(t, []error{errors.New("wat"), errors.New("lol")}, ce.getErrorSamples()) // new error added
}
