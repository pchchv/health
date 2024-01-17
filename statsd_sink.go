package health

import "bytes"

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
