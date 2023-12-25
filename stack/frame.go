package stack

// Frame contains all necessary information about to generate a line in a callstack.
type Frame struct {
	File            string
	LineNumber      int
	Name            string
	Package         string
	IsSystemPackage bool
	ProgramCounter  uintptr
}
