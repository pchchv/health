package health

import (
	"bytes"
	"net"
	"time"
)

const (
	statsdCmdKindStop statsdCmdKind = iota
	statsdCmdKindEvent
	statsdCmdKindGauge
	statsdCmdKindFlush
	statsdCmdKindDrain
	statsdCmdKindTiming
	statsdCmdKindComplete
	statsdCmdKindEventErr
)

var defaultStatsDOptions = StatsDSinkOptions{SanitizationFunc: sanitizeKey}

type StatsDSinkSanitizationFunc func(*bytes.Buffer, string)

type statsdCmdKind int

type eventKey struct {
	job    string
	event  string
	suffix string
}

type prefixBuffer struct {
	*bytes.Buffer
	prefixLen int
}

type StatsDSinkOptions struct {
	// Prefix is something like "metroid"
	// Events emitted to StatsD would be metroid.myevent.wat
	// Eg, don't include a trailing dot in the prefix.
	// It can be "", that's fine.
	Prefix string
	// SanitizationFunc sanitizes jobs and events before sending them to statsd
	SanitizationFunc StatsDSinkSanitizationFunc
	// SkipNestedEvents will skip {events,timers,gauges} from sending the job.event version
	// and will only send the event version.
	SkipNestedEvents bool
	// SkipTopLevelEvents will skip {events,timers,gauges} from sending the event version
	// and will only send the job.event version.
	SkipTopLevelEvents bool
}

type statsdEmitCmd struct {
	Kind   statsdCmdKind
	Job    string
	Event  string
	Nanos  int64
	Value  float64
	Status CompletionStatus
}

type StatsDSink struct {
	options       StatsDSinkOptions
	cmdChan       chan statsdEmitCmd
	drainDoneChan chan struct{}
	stopDoneChan  chan struct{}
	flushPeriod   time.Duration
	udpBuf        bytes.Buffer
	timingBuf     []byte
	udpConn       *net.UDPConn
	udpAddr       *net.UDPAddr
	// map of {job,event,suffix} to a re-usable buffer prefixed with the key.
	// Since each timing/gauge has a unique component (the time), we'll truncate to the prefix, write the timing,
	// and write the statsD suffix (eg, "|ms\n"). Then copy that to the UDP buffer.
	prefixBuffers map[eventKey]prefixBuffer
}

func (s *StatsDSink) flush() {
	if s.udpBuf.Len() > 0 {
		s.udpConn.WriteToUDP(s.udpBuf.Bytes(), s.udpAddr)
		s.udpBuf.Truncate(0)
	}
}

func (s *StatsDSink) writeSanitizedKeys(b *bytes.Buffer, keys ...string) {
	needDot := false
	for _, k := range keys {
		if k != "" {
			if needDot {
				b.WriteByte('.')
			}
			s.options.SanitizationFunc(b, k)
			needDot = true
		}
	}
}

func sanitizeKey(b *bytes.Buffer, s string) {
	b.Grow(len(s) + 1)
	for i := 0; i < len(s); i++ {
		si := s[i]
		if ('A' <= si && si <= 'Z') || ('a' <= si && si <= 'z') || ('0' <= si && s[i] <= '9') || si == '_' || si == '.' {
			b.WriteByte(si)
		} else {
			b.WriteByte('$')
		}
	}
}
