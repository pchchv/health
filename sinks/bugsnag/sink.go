package bugsnag

import "github.com/pchchv/health"

type cmdEventErr struct {
	Job   string
	Event string
	Err   *health.UnmutedError
	Kvs   map[string]string
}
