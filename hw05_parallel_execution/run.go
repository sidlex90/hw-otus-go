package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tCh := make(chan Task)
	var errCount int64
	var isErrorLimitExceeded bool
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tCh {
				if err := task(); err != nil {
					atomic.AddInt64(&errCount, 1)
				}
			}
		}()
	}

	maxErr := int64(m)
	for _, task := range tasks {
		if atomic.LoadInt64(&errCount) >= maxErr {
			isErrorLimitExceeded = true
		}
		if isErrorLimitExceeded {
			break
		}
		tCh <- task
	}

	close(tCh)
	wg.Wait()

	if isErrorLimitExceeded {
		return ErrErrorsLimitExceeded
	}
	return nil
}
