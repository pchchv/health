package health

type testSink struct {
	LastEmitKind   string // "Event", "EventErr", ..., "Complete"
	LastJob        string
	LastEvent      string
	LastErr        error
	LastErrEmitted bool
	LastErrUnmuted bool
	LastErrMuted   bool
	LastErrRaw     bool
	LastNanos      int64
	LastValue      float64
	LastKvs        map[string]string
	LastStatus     CompletionStatus
}

func (s *testSink) EmitEvent(job string, event string, kvs map[string]string) {
	s.LastEmitKind = "Event"
	s.LastJob = job
	s.LastEvent = event
	s.LastKvs = kvs
}

func (s *testSink) EmitEventErr(job string, event string, inputErr error, kvs map[string]string) {
	s.LastEmitKind = "EventErr"
	s.LastJob = job
	s.LastEvent = event
	s.LastKvs = kvs
	s.LastErr = inputErr

	switch inputErr := inputErr.(type) {
	case *UnmutedError:
		s.LastErrUnmuted = true
		s.LastErrEmitted = inputErr.Emitted
	case *MutedError:
		s.LastErrMuted = true
	default: // eg, case error:
		s.LastErrRaw = true
	}
}
