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
