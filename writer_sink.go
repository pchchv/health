package health

import (
	"bytes"
	"io"
	"sort"
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

func writeMapConsistently(b *bytes.Buffer, kvs map[string]string) {
	if kvs == nil {
		return
	}

	keys := make([]string, 0, len(kvs))
	for k := range kvs {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	keysLenMinusOne := len(keys) - 1

	b.WriteString(" kvs:[")
	for i, k := range keys {
		b.WriteString(k)
		b.WriteRune(':')
		b.WriteString(kvs[k])

		if i != keysLenMinusOne {
			b.WriteRune(' ')
		}
	}
	b.WriteRune(']')
}
