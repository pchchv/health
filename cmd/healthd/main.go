package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Get inputs.
	// Read from env variables for now (command line options?)
	monitoredHostPorts := getMonitoredHostPorts()
	serverHostPort := getServerHostPort()
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

func getServerHostPort() string {
	ret := os.Getenv("HEALTHD_SERVER_HOSTPORT")
	if ret == "" {
		ret = ":5031"
	}
	return ret
}
