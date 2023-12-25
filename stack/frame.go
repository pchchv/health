package stack

import (
	"runtime"
	"strings"
)

// Frame contains all necessary information about to generate a line in a callstack.
type Frame struct {
	File            string
	LineNumber      int
	Name            string
	Package         string
	IsSystemPackage bool
	ProgramCounter  uintptr
}

func packageAndName(fn *runtime.Func) (pkg string, name string) {
	name = fn.Name()
	// first remove the path prefix if there is one
	if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
		pkg += name[:lastslash] + "/"
		name = name[lastslash+1:]
	}

	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}
	return
}
