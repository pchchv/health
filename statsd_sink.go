package health

import "bytes"

type StatsDSinkSanitizationFunc func(*bytes.Buffer, string)

type eventKey struct {
	job    string
	event  string
	suffix string
}
