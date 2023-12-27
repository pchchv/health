package health

import "time"

const (
	cmdKindEvent cmdKind = iota
	cmdKindEventErr
	cmdKindTiming
	cmdKindGauge
	cmdKindComplete
)

type cmdKind int

type emitCmd struct {
	Kind   cmdKind
	Job    string
	Event  string
	Err    error
	Nanos  int64
	Value  float64
	Status CompletionStatus
}

type JsonPollingSink struct {
	intervalDuration  time.Duration
	cmdChan           chan *emitCmd
	doneChan          chan int
	doneDoneChan      chan int
	intervalsChanChan chan chan []*IntervalAggregation
}
