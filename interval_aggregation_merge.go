package health

func (intoTa *TimerAggregation) merge(fromTa *TimerAggregation) {
	intoTa.Count += fromTa.Count
	intoTa.NanosSum += fromTa.NanosSum
	intoTa.NanosSumSquares += fromTa.NanosSumSquares
	if fromTa.NanosMin < intoTa.NanosMin {
		intoTa.NanosMin = fromTa.NanosMin
	}
	if fromTa.NanosMax > intoTa.NanosMax {
		intoTa.NanosMax = fromTa.NanosMax
	}
}
