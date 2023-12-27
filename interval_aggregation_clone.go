package health

func (ec *ErrorCounter) Clone() *ErrorCounter {
	var dup = *ec
	return &dup
