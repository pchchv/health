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

func advanceNowMock(dur time.Duration) {
	nowMut.Lock()
	defer nowMut.Unlock()
	if nowMock.IsZero() {
		panic("nowMock is not set")
	}
	nowMock = nowMock.Add(dur)
}
