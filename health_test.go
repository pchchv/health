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
