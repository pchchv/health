package health

import "time"

const (
	Success CompletionStatus = iota
	ValidationError
	Panic
	Error
	Junk
)

var completionStatusToString = map[CompletionStatus]string{
	Success:         "success",
	ValidationError: "validation_error",
	Panic:           "panic",
	Error:           "error",
	Junk:            "junk",
}

// This is primarily used as syntactic sugar for libs outside this app for passing in maps easily.
// We don't rely on it internally b/c I don't want to tie interfaces to the 'health' package.
type Kvs map[string]string

type EventReceiver interface {
	Event(eventName string)
	EventKv(eventName string, kvs map[string]string)
	EventErr(eventName string, err error) error
	EventErrKv(eventName string, err error, kvs map[string]string) error
	Timing(eventName string, nanoseconds int64)
	TimingKv(eventName string, nanoseconds int64, kvs map[string]string)
	Gauge(eventName string, value float64)
	GaugeKv(eventName string, value float64, kvs map[string]string)
}

type Sink interface {
	EmitEvent(job string, event string, kvs map[string]string)
	EmitEventErr(job string, event string, err error, kvs map[string]string)
	EmitTiming(job string, event string, nanoseconds int64, kvs map[string]string)
	EmitComplete(job string, status CompletionStatus, nanoseconds int64, kvs map[string]string)
	EmitGauge(job string, event string, value float64, kvs map[string]string)
}

type CompletionStatus int

func (cs CompletionStatus) String() string {
	return completionStatusToString[cs]
}

type Job struct {
	Stream    *Stream
	JobName   string
	KeyValues map[string]string
	Start     time.Time
}

type Stream struct {
	Sinks     []Sink
	KeyValues map[string]string
	*Job
}
