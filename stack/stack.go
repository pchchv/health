package stack

type Trace struct {
	stack  []uintptr
	frames []Frame
}
