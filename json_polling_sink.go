package health

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
