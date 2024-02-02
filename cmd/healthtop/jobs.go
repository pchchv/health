package main

import (
	"fmt"
	"strings"

	"github.com/buger/goterm"
	"github.com/pchchv/health/healthd"
)

type jobOptions struct {
	Sort string
	Name string
}

func normal(text string) string {
	return fmt.Sprintf("\033[0m%s\033[0m", text)
}

func format(text string, color int, isBold bool) string {
	if isBold {
		return goterm.Bold(goterm.Color(text, color))
	} else {
		return normal(goterm.Color(text, color))
	}
}

func printCellString(text string, table *goterm.Table, isBold, isGreen, isRed bool) {
	color := goterm.BLACK
	if isGreen {
		color = goterm.GREEN
	} else if isRed {
		color = goterm.RED
	}
	fmt.Fprintf(table, "%s\t", format(text, color, isBold))
}

func printCellInt64(val int64, table *goterm.Table, isBold, isGreen, isRed bool) {
	printCellString(fmt.Sprint(val), table, isBold, isGreen, isRed)
}

func printCellNanos(nanos int64, table *goterm.Table, isBold, isGreen, isRed bool) {
	var units string
	switch {
	case nanos > 2000000:
		units = "ms"
		nanos /= 1000000
	case nanos > 1000:
		units = "μs"
		nanos /= 1000
	default:
		units = "ns"
	}
	printCellString(fmt.Sprintf("%d %s", nanos, units), table, isBold, isGreen, isRed)
}

func printJob(table *goterm.Table, job *healthd.Job) {
	fullSuccess := job.Count == job.CountSuccess
	printCellString(job.Name, table, true, false, false)
	printCellInt64(job.Count, table, false, fullSuccess, false)
	printCellInt64(job.CountSuccess, table, fullSuccess, fullSuccess, false)
	printCellInt64(job.CountValidationError, table, job.CountValidationError > 0, false, job.CountValidationError > 0)
	printCellInt64(job.CountPanic, table, job.CountPanic > 0, false, job.CountPanic > 0)
	printCellInt64(job.CountError, table, job.CountError > 0, false, job.CountError > 0)
	printCellInt64(job.CountJunk, table, job.CountJunk > 0, false, job.CountJunk > 0)
	printCellNanos(int64(job.NanosAvg), table, true, false, false)
	printCellNanos(int64(job.NanosStdDev), table, false, false, false)
	printCellNanos(job.NanosMin, table, false, false, false)
	printCellNanos(job.NanosMax, table, false, false, false)
	printCellNanos(job.NanosSum, table, false, false, false)
	fmt.Fprintf(table, "\n")
}

func printJobs(lastApiResponse *healthd.ApiResponseJobs, status *healthdStatus) {
	goterm.Clear() // clear current screen
	goterm.MoveCursor(1, 1)
	defer goterm.Flush()
	goterm.Println("Current Time:", status.FmtNow(), "   Status:", status.FmtStatus())

	if lastApiResponse == nil {
		goterm.Println("no data yet")
		return
	}

	columns := []string{
		"Job",
		"Total Count",
		"Success",
		"ValidationError",
		"Panic",
		"Error",
		"Junk",
		"Avg Response Time",
		"Stddev",
		"Min",
		"Max",
		"Total",
	}
	for i, s := range columns {
		columns[i] = goterm.Bold(goterm.Color(s, goterm.BLACK))
	}

	table := goterm.NewTable(0, goterm.Width()-1, 5, ' ', 0)
	fmt.Fprintf(table, "%s\n", strings.Join(columns, "\t"))
	for _, job := range lastApiResponse.Jobs {
		printJob(table, job)
	}
	goterm.Println(table)
}
