package hw05parallelexecution

import (
	"fmt"
)

// worker выполняет задачи
type Worker struct {
	id             int
	WorkersChannel chan chan Task // когда воркер свободен - он добавит свой канал сюда
	WorkChannel    chan Task      // канал с задачей
	errorChannel   chan error     // канал с ошибками
	end            chan bool      // канал для остановки воркера
}

func (w *Worker) Start() {

	fmt.Printf("W-%v started...\n", w.id)
	go func() {
		for {
			w.WorkersChannel <- w.WorkChannel // помещаем канал для задачи в общий канал
			select {
			case task, ok := <-w.WorkChannel:
				if !ok {
					w.Stop()
				}
				fmt.Printf("🛑 W-%v working on task...\n", w.id) // получили задачу
				err := task()
				w.errorChannel <- err
				fmt.Printf("✔️ W-%v task done\n", w.id)

			case <-w.end:
				return

			}

		}
	}()
}

func (w *Worker) Stop() {
	fmt.Printf("⛔ W-%v stopping\n", w.id)
	w.end <- true
}
