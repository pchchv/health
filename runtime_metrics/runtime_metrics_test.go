package runtime_metrics

type testReceiver struct {
	gauges map[string]float64
}

func (t *testReceiver) Event(eventName string) {}

func (t *testReceiver) EventKv(eventName string, kvs map[string]string) {}

func (t *testReceiver) EventErr(eventName string, err error) error {
	return nil
}

func (t *testReceiver) EventErrKv(eventName string, err error, kvs map[string]string) error {
	return nil
}

func (t *testReceiver) Timing(eventName string, nanoseconds int64) {}

func (t *testReceiver) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {}
