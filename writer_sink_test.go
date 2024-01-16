package health

import (
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testErr             = errors.New("test error")
	basicEventRegexp    = regexp.MustCompile("\\[[^\\]]+\\]: job:(.+) event:(.+)")
	basicEventErrRegexp = regexp.MustCompile("\\[[^\\]]+\\]: job:(.+) event:(.+) err:(.+)")
	basicTimingRegexp   = regexp.MustCompile("\\[[^\\]]+\\]: job:(.+) event:(.+) time:(.+)")
	kvsEventRegexp      = regexp.MustCompile("\\[[^\\]]+\\]: job:(.+) event:(.+) kvs:\\[(.+)\\]")
	kvsEventErrRegexp   = regexp.MustCompile("\\[[^\\]]+\\]: job:(.+) event:(.+) err:(.+) kvs:\\[(.+)\\]")
)

func TestWriterSinkEmitEventBasic(t *testing.T) {
	var b bytes.Buffer
	sink := WriterSink{&b}
	sink.EmitEvent("myjob", "myevent", nil)

	str := b.String()

	result := basicEventRegexp.FindStringSubmatch(str)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, "myjob", result[1])
	assert.Equal(t, "myevent", result[2])
}

func TestWriterSinkEmitEventKvs(t *testing.T) {
	var b bytes.Buffer
	sink := WriterSink{&b}
	sink.EmitEvent("myjob", "myevent", map[string]string{"wat": "ok", "another": "thing"})

	str := b.String()

	result := kvsEventRegexp.FindStringSubmatch(str)
	assert.Equal(t, 4, len(result))
	assert.Equal(t, "myjob", result[1])
	assert.Equal(t, "myevent", result[2])
	assert.Equal(t, "another:thing wat:ok", result[3])
}

func TestWriterSinkEmitEventErrBasic(t *testing.T) {
	var b bytes.Buffer
	sink := WriterSink{&b}
	sink.EmitEventErr("myjob", "myevent", testErr, nil)

	str := b.String()

	result := basicEventErrRegexp.FindStringSubmatch(str)
	assert.Equal(t, 4, len(result))
	assert.Equal(t, "myjob", result[1])
	assert.Equal(t, "myevent", result[2])
	assert.Equal(t, testErr.Error(), result[3])
}

func TestWriterSinkEmitEventErrKvs(t *testing.T) {
	var b bytes.Buffer
	sink := WriterSink{&b}
	sink.EmitEventErr("myjob", "myevent", testErr, map[string]string{"wat": "ok", "another": "thing"})

	str := b.String()

	result := kvsEventErrRegexp.FindStringSubmatch(str)
	assert.Equal(t, 5, len(result))
	assert.Equal(t, "myjob", result[1])
	assert.Equal(t, "myevent", result[2])
	assert.Equal(t, testErr.Error(), result[3])
	assert.Equal(t, "another:thing wat:ok", result[4])
}

func TestWriterSinkEmitTimingBasic(t *testing.T) {
	var b bytes.Buffer
	sink := WriterSink{&b}
	sink.EmitTiming("myjob", "myevent", 1204000, nil)

	str := b.String()

	result := basicTimingRegexp.FindStringSubmatch(str)
	assert.Equal(t, 4, len(result))
	assert.Equal(t, "myjob", result[1])
	assert.Equal(t, "myevent", result[2])
	assert.Equal(t, "1204 μs", result[3])
}
