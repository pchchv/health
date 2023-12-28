package health

import "testing"

// Let's leverage clone's fixture data and make sure if
// it is possible to merge a new blank aggregation to get the same data.
func TestMergeBasic(t *testing.T) {
	setNowMock("2011-09-09T23:36:13Z")
	defer resetNowMock()

	a := aggregatorWithData()
	intAgg := a.intervalAggregations[0]
	assertAggregationData(t, intAgg)
	newAgg := NewIntervalAggregation(intAgg.IntervalStart)
	newAgg.Merge(intAgg)
	assertAggregationData(t, newAgg)
}
