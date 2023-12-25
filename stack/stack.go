package stack

import "runtime"

var MaxStackDepth = 50 // maximum number of stackframes on any error

type Trace struct {
	stack  []uintptr
	frames []Frame
}

func NewTrace(skip int) *Trace {
	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2+skip, stack)
	return &Trace{
		stack: stack[:length],
	}
}
