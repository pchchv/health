package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	monitoredHostPorts := getMonitoredHostPorts()
}

func getMonitoredHostPorts() []string {
	hps := os.Getenv("HEALTHD_MONITORED_HOSTPORTS")
	if hps == "" {
		fmt.Println("no hosts to monitor. Pass them in with the environment variable HEALTHD_MONITORED_HOSTPORTS")
		fmt.Println("example: $ HEALTHD_MONITORED_HOSTPORTS=web31:5020,web32:5020 ./healthd")
		os.Exit(1)
	}
	return strings.Split(hps, ",")
}
