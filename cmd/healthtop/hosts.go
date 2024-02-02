package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/buger/goterm"
	"github.com/pchchv/health/healthd"
)

func printHost(table *goterm.Table, host *healthd.HostStatus) {
	var status string
	success := host.LastCode == 200 && host.LastErr == ""
	if success {
		status = "Success"
	} else if host.LastCheckTime.IsZero() {
		status = "Unknown"
	} else {
		status = "Failure: " + host.LastErr
	}

	printCellString(host.HostPort, table, true, false, false)
	printCellString(status, table, false, success, !success)
	printCellString(host.LastCheckTime.Format(time.RFC1123), table, false, false, false)
	printCellNanos(int64(host.LastNanos), table, false, false, false)
	fmt.Fprintf(table, "\n")
}

func printHosts(lastApiResponse *healthd.ApiResponseHosts, status *healthdStatus) {
	goterm.Clear() // clear current screen
	goterm.MoveCursor(1, 1)
	defer goterm.Flush()
	goterm.Println("Current Time:", status.FmtNow(), "   Status:", status.FmtStatus())

	if lastApiResponse == nil {
		goterm.Println("no data yet")
		return
	}

	columns := []string{
		"Host:Port",
		"Status",
		"Last Checked",
		"Last Response Time",
	}
	for i, s := range columns {
		columns[i] = goterm.Bold(goterm.Color(s, goterm.BLACK))
	}

	table := goterm.NewTable(0, goterm.Width()-1, 5, ' ', 0)
	fmt.Fprintf(table, "%s\n", strings.Join(columns, "\t"))
	for _, host := range lastApiResponse.Hosts {
		printHost(table, host)
	}

	goterm.Println(table)
}
