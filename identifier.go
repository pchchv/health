package health

import (
	"fmt"
	"os"
)

var Identifier = getIdentifier()

func getIdentifier() string {
	pid := os.Getpid()
	host, err := os.Hostname()
	if err != nil {
		host = "hostname errored"
	}

	return fmt.Sprintf("%s.%d", host, pid)
}
