package health

type MutedError struct {
	Err error
}

func (e *MutedError) Error() string {
	return e.Err.Error()
}
