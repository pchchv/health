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

func NewJsonPollingSink(intervalDuration time.Duration, retain time.Duration) *JsonPollingSink {
	const buffSize = 4096 // random-ass-guess
	s := &JsonPollingSink{
		intervalDuration:  intervalDuration,
		cmdChan:           make(chan *emitCmd, buffSize),
		doneChan:          make(chan int),
		doneDoneChan:      make(chan int),
		intervalsChanChan: make(chan chan []*IntervalAggregation),
	}

	go startAggregator(intervalDuration, retain, s)

	return s
}

func (s *JsonPollingSink) ShutdownServer() {
	s.doneChan <- 1
	<-s.doneDoneChan
}
