package main

import (
	"fmt"
	"time"
)

type healthdStatus struct {
	lastSuccessAt time.Time
	lastErrorAt   time.Time
	lastError     error
}

func (s *healthdStatus) FmtNow() string {
	return time.Now().Format(time.RFC1123)
}

func (s *healthdStatus) FmtStatus() string {
	if s.lastErrorAt.IsZero() && s.lastSuccessAt.IsZero() {
		return "[starting...]"
	} else if s.lastErrorAt.After(s.lastSuccessAt) {
		return fmt.Sprint("[error: '", s.lastError.Error(), "'    LastErrorAt: ", s.lastErrorAt.Format(time.RFC1123), "]")
	} else {
		return "[success]"
	}
}

func main() {}
