package health

import "io"

type JsonWriterSink struct {
	io.Writer
}
