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

type Worker struct {
	TasksChannel chan Task
}

func (w *Worker) Start(wg *sync.WaitGroup, sum *ErrSumm) {
	go func() {
		defer wg.Done()
		for {
			task, ok := <-w.TasksChannel
			if !ok {
				return
			}

			res := task()
			if res != nil {
				sum.Lock()
				sum.totalErrors++
				sum.Unlock()
			}
		}
	}()
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	errCounter := ErrSumm{}
	tasksChannel := make(chan Task)

	// Запуск воркеров
	for w := 0; w < n; w++ {
		worker := Worker{
			TasksChannel: tasksChannel,
		}
		wg.Add(1)
		worker.Start(&wg, &errCounter)
	}

	// Отправка заданий в канал
	var totalErrors int
	for i := 0; i < len(tasks); i++ {
		errCounter.Lock()
		totalErrors = errCounter.totalErrors
		errCounter.Unlock()
		if totalErrors > 0 && totalErrors >= m {
			break
		}
		tasksChannel <- tasks[i]
	}

	// Закрываем канал и ждем остановки горутин
	close(tasksChannel)
	wg.Wait()

	if totalErrors > 0 && totalErrors >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
