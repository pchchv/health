package health

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var basicEventRegexp = regexp.MustCompile("\\[[^\\]]+\\]: job:(.+) event:(.+)")

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
