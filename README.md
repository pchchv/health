# health [![GoDoc](https://godoc.org/github.com/pchchv/health?status.png)](https://godoc.org/github.com/pchchv/health)

health allows you to instrument your service for logging and metrics, and then send that to log files, StatsD, Bugsnag, or to be polled and aggregate via JSON API.

health also ships with a New Relic-like aggregator (called healthd) that shows the slowest endpoints, biggest error producers, highest throughput endpoints, and so on.
