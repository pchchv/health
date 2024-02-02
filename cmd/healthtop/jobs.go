package main

import (
	"fmt"

	"github.com/buger/goterm"
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
