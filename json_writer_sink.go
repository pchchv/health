package health

import (
	"encoding/json"
	"io"
)

type JsonWriterSink struct {
	io.Writer
}

func (j *JsonWriterSink) EmitEvent(job string, event string, kvs map[string]string) {
	b, err := json.Marshal(struct {
		Job       string
		Event     string
		Timestamp string
		Kvs       map[string]string
	}{job, event, timestamp(), kvs})

	if err != nil {
		return
	}
	j.Write(b)
}

func (j *JsonWriterSink) EmitEventErr(job string, event string, err error, kvs map[string]string) {
	b, err := json.Marshal(struct {
		Job       string
		Event     string
		Timestamp string
		Err       string
		Kvs       map[string]string
	}{job, event, timestamp(), err.Error(), kvs})

	if err != nil {
		return
	}
	j.Write(b)
}
