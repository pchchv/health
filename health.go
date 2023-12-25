package health

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
