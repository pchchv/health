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

// Frames returns an array of frames containing information about the stack.
func (t *Trace) Frames() []Frame {
	if t.frames == nil {
		t.frames = make([]Frame, len(t.stack))
		for i, pc := range t.stack {
			t.frames[i] = NewFrame(pc)
		}
	}
	return t.frames
}
