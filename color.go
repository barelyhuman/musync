package main

import (
	"runtime"
	"strings"
)

var (
	terminalGreen = "\033[32m"
	terminalRed   = "\033[31m"
	terminalReset = "\033[0m"
)

func colorString(color string, toWrite string) string {
	if runtime.GOOS == "windows" {
		terminalGreen = ""
		terminalRed = ""
		terminalReset = ""
	}

	var sb strings.Builder
	sb.WriteString(color)
	sb.WriteString(toWrite)
	sb.WriteString(terminalReset)
	return sb.String()
}
