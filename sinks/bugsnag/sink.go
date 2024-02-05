package bugsnag

import (
	"fmt"
	"os"

	"github.com/pchchv/health"
)

type cmdEventErr struct {
	Job   string
	Event string
	Err   *health.UnmutedError
	Kvs   map[string]string
}

// This sink emits to a StatsD deaemon by sending it a UDP packet.
type Sink struct {
	*Config
	cmdChan  chan *cmdEventErr
	doneChan chan int
}

func errorProcessingLoop(sink *Sink) {
	cmdChan := sink.cmdChan
	doneChan := sink.doneChan

PROCESSING_LOOP:
	for {
		select {
		case <-doneChan:
			break PROCESSING_LOOP
		case cmd := <-cmdChan:
			if err := Notify(sink.Config, cmd.Job, cmd.Event, cmd.Err, cmd.Err.Stack, cmd.Kvs); err != nil {
				fmt.Fprintf(os.Stderr, "bugsnag.Notify: could not notify bugsnag. err=%v\n", err)
			}
		}
	}
}
