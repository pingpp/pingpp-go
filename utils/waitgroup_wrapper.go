package utils

import (
	"sync"
	"time"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (wg *WaitGroupWrapper) Wrap(cb func()) {
	wg.Add(1)
	go func() {
		cb()
		wg.Done()
	}()
}

func (wg *WaitGroupWrapper) WaitTimeout(timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
