package hw05parallelexecution

import (
	"errors"
	"fmt"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {

	errorsChannel := make(chan error, 1000)    // результат выполнения таска
	queueWorkers := make(chan chan Task, 1000) // канал, куда воркеры помещают свой канал, когда освободятся
	workersPool := make([]*Worker, n)          // пул воркеров

	// create dispatcher with n workers
	dispatcher := Dispatcher{
		workersPool:   workersPool,
		TaskChan:      make(chan Task),
		QueueWorkers:  queueWorkers,
		ErrorsChannel: errorsChannel,
		End:           make(chan bool),
	}

	// Start dispatcher processor and workers.
	dispatcher.StartWorkers()

	var errorsLimitReached bool

	// go func() {
	// 	// adding tasks to dispatcher
	// 	for _, task := range tasks {
	// 		if errorsLimitReached {
	// 			return
	// 		}
	// 		dispatcher.AddTaskToProcess(task)
	// 	}

	// }()

	// // count errrors
	// var errTotal int
	// for i := 0; i < len(tasks); i++ {
	// 	res := <-errorsChannel
	// 	if res != nil {
	// 		errTotal++
	// 		if errTotal >= m {
	// 			errorsLimitReached = true
	// 			break
	// 		}
	// 	}
	// }

	for it, task := range tasks {
		if dispatcher.ErrorsCount >= m {
			errorsLimitReached = true
			break
		}
		fmt.Println("Sending task", it)
		dispatcher.AddTaskToProcess(task)
	}

	dispatcher.End <- true
	dispatcher.StopWorkers()
	if errorsLimitReached {
		return ErrErrorsLimitExceeded
	}
	return nil
}
