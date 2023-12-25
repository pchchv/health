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

func (j *Job) KeyValue(key string, value string) *Job {
	if j.KeyValues == nil {
		j.KeyValues = make(map[string]string)
	}
	j.KeyValues[key] = value
	return j
}

func (j *Job) Event(eventName string) {
	allKvs := j.mergedKeyValues(nil)
	for _, sink := range j.Stream.Sinks {
		sink.EmitEvent(j.JobName, eventName, allKvs)
	}
}

func (j *Job) mergedKeyValues(instanceKvs map[string]string) map[string]string {
	var allKvs map[string]string

	// Count how many maps actually have contents in them.
	// If it's 0 or 1, we won't allocate a new map.
	// Also, optimistically set allKvs.
	// We might use it or we might overwrite the value with a newly made map.
	var kvCount = 0
	if len(j.KeyValues) > 0 {
		kvCount += 1
		allKvs = j.KeyValues
	}

	if len(j.Stream.KeyValues) > 0 {
		kvCount += 1
		allKvs = j.Stream.KeyValues
	}

	if len(instanceKvs) > 0 {
		kvCount += 1
		allKvs = instanceKvs
	}

	if kvCount > 1 {
		allKvs = make(map[string]string)
		for k, v := range j.Stream.KeyValues {
			allKvs[k] = v
		}
		for k, v := range j.KeyValues {
			allKvs[k] = v
		}
		for k, v := range instanceKvs {
			allKvs[k] = v
		}
	}
	return allKvs
}

type Stream struct {
	Sinks     []Sink
	KeyValues map[string]string
	*Job
}

func NewStream() *Stream {
	s := &Stream{}
	s.Job = s.NewJob("general")
	return s
}

func (s *Stream) NewJob(name string) *Job {
	return &Job{
		Stream:  s,
		JobName: name,
		Start:   time.Now(),
	}
}

func (s *Stream) AddSink(sink Sink) *Stream {
	s.Sinks = append(s.Sinks, sink)
	return s
}

func (s *Stream) KeyValue(key string, value string) *Stream {
	if s.KeyValues == nil {
		s.KeyValues = make(map[string]string)
	}
	s.KeyValues[key] = value
	return s
}
