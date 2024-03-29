package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pchchv/health"
	"github.com/pchchv/health/healthd"
)

func main() {
	// Get inputs.
	// Read from env variables for now (command line options?)
	monitoredHostPorts := getMonitoredHostPorts()
	serverHostPort := getServerHostPort()
	healthHostPort := getHealthHostPort()

	// Monitor ourselves.
	// This will make our own instrumentation show up in the healthd output.
	monitoredHostPorts = append(monitoredHostPorts, healthHostPort)

	// Setup our health stream.
	// Log to stdout and a setup an polling sink
	stream := health.NewStream()
	stream.AddSink(&health.WriterSink{Writer: os.Stdout})
	jsonPollingSink := health.NewJsonPollingSink(time.Minute, time.Minute*5)
	jsonPollingSink.StartServer(healthHostPort)
	stream.AddSink(jsonPollingSink)

	// Say we're starting!
	stream.EventKv("starting", health.Kvs{
		"monitored_host_ports": strings.Join(monitoredHostPorts, ","),
		"server_host_port":     serverHostPort,
		"health_host_port":     healthHostPort,
	})

	// Start the healthd aggregators in a goroutine(s)
	healthd.StartNewHealthD(monitoredHostPorts, serverHostPort, stream)

	// Block
	select {}
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

func getHealthHostPort() string {
	ret := os.Getenv("HEALTH_HOSTPORT")
	if ret == "" {
		ret = ":5030"
	}
	return ret
}
