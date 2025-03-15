package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tCh := make(chan Task)
	var errCount int
	var isErrorLimitExceeded bool
	wg, mu := sync.WaitGroup{}, sync.RWMutex{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(tCh)
		for _, task := range tasks {
			mu.RLock()
			if errCount >= m {
				isErrorLimitExceeded = true
			}
			mu.RUnlock()
			if isErrorLimitExceeded {
				break
			}
			tCh <- task
		}
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tCh {
				if err := task(); err != nil {
					mu.Lock()
					errCount++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	if isErrorLimitExceeded {
		return ErrErrorsLimitExceeded
	}
	return nil
}
