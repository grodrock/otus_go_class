package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrSumm struct {
	sync.Mutex
	totalErrors int
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	errCounter := ErrSumm{}

	for i := 0; i < len(tasks); {
		for w := 0; w < n && i < len(tasks); w++ {
			wg.Add(1)
			go func(t Task) {
				defer wg.Done()
				res := t()
				if res != nil {
					errCounter.Lock()
					errCounter.totalErrors++
					errCounter.Unlock()
				}
			}(tasks[i])
			i++
		}
		wg.Wait()
		if errCounter.totalErrors >= m {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}
