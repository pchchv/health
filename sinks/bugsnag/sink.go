package bugsnag

import "github.com/pchchv/health"

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
