package health

import "github.com/pchchv/health/stack"

type MutedError struct {
	Err error
}

func (e *MutedError) Error() string {
	return e.Err.Error()
}

type UnmutedError struct {
	Err     error
	Stack   *stack.Trace
	Emitted bool
}

func (e *UnmutedError) Error() string {
	return e.Err.Error()
}

func Mute(err error) *MutedError {
	return &MutedError{Err: err}
}
