package health

import (
	"io"
	"time"
)

// This sink writes bytes in a format that a human might like to read in a logfile
// This can be used to log to Stdout:
//
//	.AddSink(&WriterSink{os.Stdout})
//
// And to a file:
//
//	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
//	.AddSink(&WriterSink{f})
//
// And to syslog:
//
//	w, err := syslog.New(LOG_INFO, "wat")
//	.AddSink(&WriterSink{w})
type WriterSink struct {
	io.Writer
}

func timestamp() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}
