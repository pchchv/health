package healthd

import "github.com/pchchv/health"

type apiContext struct {
	hd *HealthD
	*health.Job
}
