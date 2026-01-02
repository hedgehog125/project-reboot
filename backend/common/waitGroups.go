package common

import (
	"sync"
	"sync/atomic"
)

type WaitGroupWithCounter struct {
	sync.WaitGroup
	waitingCount atomic.Int64
}

func (wg *WaitGroupWithCounter) WaitingCount() int64 {
	return wg.waitingCount.Load()
}
func (wg *WaitGroupWithCounter) Add(delta int) {
	wg.waitingCount.Add(int64(delta))
	wg.WaitGroup.Add(delta)
}
func (wg *WaitGroupWithCounter) Done() {
	wg.waitingCount.Add(-1)
	wg.WaitGroup.Done()
}
