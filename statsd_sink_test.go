package health

import (
	"fmt"
	"net"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testAddr = "127.0.0.1:7890"

func TestStatsDSinkPeriodicPurge(t *testing.T) {
	sink, err := NewStatsDSink(testAddr, &StatsDSinkOptions{Prefix: "metroid"})
	assert.NoError(t, err)

	// Stop the sink, set a smaller flush period, and start it agian
	sink.Stop()
	sink.flushPeriod = 1 * time.Millisecond
	go sink.loop()
	defer sink.Stop()

	listenFor(t, []string{"metroid.my.event:1|c\nmetroid.my.job.my.event:1|c\n"}, func() {
		sink.EmitEvent("my.job", "my.event", nil)
		time.Sleep(10 * time.Millisecond)
	})
}

func TestStatsDSinkPacketLimit(t *testing.T) {
	sink, err := NewStatsDSink(testAddr, &StatsDSinkOptions{Prefix: "metroid", SkipNestedEvents: true})
	assert.NoError(t, err)

	// s is 101 bytes
	s := "metroid." + strings.Repeat("a", 88) + ":1|c\n"

	// expect 1 packet that is 14*101=1414 bytes, and the next one to be 101 bytes
	listenFor(t, []string{strings.Repeat(s, 14), s}, func() {
		for i := 0; i < 15; i++ {
			sink.EmitEvent("my.job", strings.Repeat("a", 88), nil)
		}

		sink.Drain()
	})
}

func TestStatsDSinkEmitEventPrefix(t *testing.T) {
	sink, err := NewStatsDSink(testAddr, &StatsDSinkOptions{Prefix: "metroid"})
	defer sink.Stop()
	assert.NoError(t, err)
	listenFor(t, []string{"metroid.my.event:1|c\nmetroid.my.job.my.event:1|c\n"}, func() {
		sink.EmitEvent("my.job", "my.event", nil)
		sink.Drain()
	})
}

func TestStatsDSinkEmitEventNoPrefix(t *testing.T) {
	sink, err := NewStatsDSink(testAddr, nil)
	defer sink.Stop()
	assert.NoError(t, err)
	listenFor(t, []string{"my.event:1|c\nmy.job.my.event:1|c\n"}, func() {
		sink.EmitEvent("my.job", "my.event", nil)
		sink.Drain()
	})
}

func listenFor(t *testing.T, msgs []string, f func()) {
	c, err := net.ListenPacket("udp", testAddr)
	defer c.Close()
	assert.NoError(t, err)

	f()

	buf := make([]byte, 10000)
	for _, msg := range msgs {
		err = c.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
		assert.NoError(t, err)
		nbytes, _, err := c.ReadFrom(buf)
		assert.NoError(t, err)
		if err == nil {
			gotMsg := string(buf[0:nbytes])
			if gotMsg != msg {
				t.Errorf("expected UPD packet %s but got %s\n", msg, gotMsg)
			}
		}
	}
}

func callerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}
