package health

import (
	"bytes"
	"net"
	"strconv"
	"time"
)

const (
	maxUdpBytes                     = 1440 // 1500(Ethernet MTU) - 60(Max UDP header size
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

func (s *StatsDSink) getPrefixBuffer(job, event, suffix string) prefixBuffer {
	key := eventKey{job, event, suffix}
	b, ok := s.prefixBuffers[key]
	if !ok {
		b.Buffer = &bytes.Buffer{}
		s.writeSanitizedKeys(b.Buffer, s.options.Prefix, job, event, suffix)
		b.WriteByte(':')
		b.prefixLen = b.Len()
		// 123456789.99|ms\n 16 bytes. timing value represents 11 days max
		b.Grow(16)
		s.prefixBuffers[key] = b
	} else {
		b.Truncate(b.prefixLen)
	}
	return b
}

func (s *StatsDSink) writeNanosToTimingBuf(nanos int64) {
	s.timingBuf = s.timingBuf[0:0]
	if nanos >= 10e6 {
		// More than 10 milliseconds. We'll just print as an integer
		s.timingBuf = strconv.AppendInt(s.timingBuf, nanos/1e6, 10)
	} else {
		s.timingBuf = strconv.AppendFloat(s.timingBuf, float64(nanos)/float64(time.Millisecond), 'f', 2, 64)
	}
}

// assumes b is a well-formed statsd metric like "job.event:1|c\n" (including newline)
func (s *StatsDSink) writeStatsDMetric(b []byte) {
	lenb := len(b)
	if lenb == 0 {
		return
	}

	// single metric exceeds limit. sad day.
	if lenb > maxUdpBytes {
		return
	}

	lenUdpBuf := s.udpBuf.Len()

	if (lenb + lenUdpBuf) > maxUdpBytes {
		s.udpConn.WriteToUDP(s.udpBuf.Bytes(), s.udpAddr)
		s.udpBuf.Truncate(0)
	}

	s.udpBuf.Write(b)
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

func (s *StatsDSink) processEvent(job string, event string) {
	if !s.options.SkipTopLevelEvents {
		pb := s.getPrefixBuffer("", event, "")
		pb.WriteString("1|c\n")
		s.writeStatsDMetric(pb.Bytes())
	}

	if !s.options.SkipNestedEvents {
		pb := s.getPrefixBuffer(job, event, "")
		pb.WriteString("1|c\n")
		s.writeStatsDMetric(pb.Bytes())
	}
}

func (s *StatsDSink) processEventErr(job string, event string) {
	if !s.options.SkipTopLevelEvents {
		pb := s.getPrefixBuffer("", event, "error")
		pb.WriteString("1|c\n")
		s.writeStatsDMetric(pb.Bytes())
	}

	if !s.options.SkipNestedEvents {
		pb := s.getPrefixBuffer(job, event, "error")
		pb.WriteString("1|c\n")
		s.writeStatsDMetric(pb.Bytes())
	}
}

func (s *StatsDSink) processTiming(job string, event string, nanos int64) {
	s.writeNanosToTimingBuf(nanos)
	if !s.options.SkipTopLevelEvents {
		pb := s.getPrefixBuffer("", event, "")
		pb.Write(s.timingBuf)
		pb.WriteString("|ms\n")
		s.writeStatsDMetric(pb.Bytes())
	}

	if !s.options.SkipNestedEvents {
		pb := s.getPrefixBuffer(job, event, "")
		pb.Write(s.timingBuf)
		pb.WriteString("|ms\n")
		s.writeStatsDMetric(pb.Bytes())
	}
}

func (s *StatsDSink) processGauge(job string, event string, value float64) {
	s.timingBuf = s.timingBuf[0:0]
	prec := 2
	if (value < 0.1) && (value > -0.1) {
		prec = -1
	}

	s.timingBuf = strconv.AppendFloat(s.timingBuf, value, 'f', prec, 64)

	if !s.options.SkipTopLevelEvents {
		pb := s.getPrefixBuffer("", event, "")
		pb.Write(s.timingBuf)
		pb.WriteString("|g\n")
		s.writeStatsDMetric(pb.Bytes())
	}

	if !s.options.SkipNestedEvents {
		pb := s.getPrefixBuffer(job, event, "")
		pb.Write(s.timingBuf)
		pb.WriteString("|g\n")
		s.writeStatsDMetric(pb.Bytes())
	}
}

func (s *StatsDSink) processComplete(job string, status CompletionStatus, nanos int64) {
	s.writeNanosToTimingBuf(nanos)
	statusString := status.String()
	pb := s.getPrefixBuffer(job, "", statusString)
	pb.Write(s.timingBuf)
	pb.WriteString("|ms\n")
	s.writeStatsDMetric(pb.Bytes())
}

func (s *StatsDSink) processCmd(cmd *statsdEmitCmd) {
	switch cmd.Kind {
	case statsdCmdKindEvent:
		s.processEvent(cmd.Job, cmd.Event)
	case statsdCmdKindEventErr:
		s.processEventErr(cmd.Job, cmd.Event)
	case statsdCmdKindTiming:
		s.processTiming(cmd.Job, cmd.Event, cmd.Nanos)
	case statsdCmdKindGauge:
		s.processGauge(cmd.Job, cmd.Event, cmd.Value)
	case statsdCmdKindComplete:
		s.processComplete(cmd.Job, cmd.Status, cmd.Nanos)
	}
}

func (s *StatsDSink) loop() {
	cmdChan := s.cmdChan
	ticker := time.NewTicker(s.flushPeriod)
	go func() {
		for range ticker.C {
			cmdChan <- statsdEmitCmd{Kind: statsdCmdKindFlush}
		}
	}()

LOOP:
	for cmd := range cmdChan {
		switch cmd.Kind {
		case statsdCmdKindDrain:
		DRAIN_LOOP:
			for {
				select {
				case cmd := <-cmdChan:
					s.processCmd(&cmd)
				default:
					s.flush()
					s.drainDoneChan <- struct{}{}
					break DRAIN_LOOP
				}
			}
		case statsdCmdKindStop:
			s.stopDoneChan <- struct{}{}
			break LOOP
		case statsdCmdKindFlush:
			s.flush()
		default:
			s.processCmd(&cmd)
		}
	}

	ticker.Stop()
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
