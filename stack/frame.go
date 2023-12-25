package stack

import (
	"fmt"
	"runtime"
	"strings"
)

var goroot = runtime.GOROOT()

// Frame contains all necessary information about to generate a line in a callstack.
type Frame struct {
	File            string
	LineNumber      int
	Name            string
	Package         string
	IsSystemPackage bool
	ProgramCounter  uintptr
}


// Func returns the function that this stackframe corresponds to.
func (frame *Frame) Func() *runtime.Func {
	if frame.ProgramCounter == 0 {
		return nil
	}
	return runtime.FuncForPC(frame.ProgramCounter)
}

func (frame *Frame) String() string {
	return fmt.Sprintf("%s:%d %s", frame.File, frame.LineNumber, frame.Name)
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

// isSystemPackage returns true iff the package is a system package like 'runtime' or 'net/http'.
func isSystemPackage(file, pkg string) bool {
	return strings.HasPrefix(file, goroot)
}
