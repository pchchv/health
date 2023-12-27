package health

func (ec *ErrorCounter) Clone() *ErrorCounter {
	var dup = *ec
	return &dup
}

func (ta *TimerAggregation) Clone() *TimerAggregation {
	var dup = *ta
	return &dup
}

func (am *aggregationMaps) Clone() *aggregationMaps {
	dup := &aggregationMaps{}
	dup.initAggregationMaps()
	for k, v := range am.Events {
		dup.Events[k] = v
	}

	for k, v := range am.Gauges {
		dup.Gauges[k] = v
	}

	for k, v := range am.Timers {
		dup.Timers[k] = v.Clone()
	}

	for k, v := range am.EventErrs {
		dup.EventErrs[k] = v.Clone()
	}
	return dup
}
