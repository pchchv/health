package health

type testJsonEvent struct {
	Job         string
	Event       string
	Timestamp   string
	Err         string
	Nanoseconds int64
	Value       float64
	Status      string
	Kvs         map[string]string
}
