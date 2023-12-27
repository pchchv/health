package health

func (ec *ErrorCounter) Clone() *ErrorCounter {
	var dup = *ec
	return &dup
}

func (ta *TimerAggregation) Clone() *TimerAggregation {
	var dup = *ta
	return &dup
}
