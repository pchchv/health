package healthd

import (
	"sync"
	"time"
)

var nowMock time.Time
var nowMut sync.RWMutex

func now() time.Time {
	nowMut.RLock()
	defer nowMut.RUnlock()
	if nowMock.IsZero() {
		return time.Now()
	}
	return nowMock
}
