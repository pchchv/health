package runtime_metrics

type testReceiver struct {
	gauges map[string]float64
}

func (t *testReceiver) Event(eventName string) {}

func (t *testReceiver) EventKv(eventName string, kvs map[string]string) {}

func (t *testReceiver) EventErr(eventName string, err error) error {
	return nil
}
