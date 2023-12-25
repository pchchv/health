package stack

type someT struct{}

func (s someT) level2() *Trace {
	return NewTrace(0)
}

func (s someT) level1() *Trace {
	return s.level2()
}

func (s someT) level0() *Trace {
	return s.level1()
}

func level2() *Trace {
	return NewTrace(0)
}

func level1() *Trace {
	return level2()
}

func level0() *Trace {
	return level1()
}
