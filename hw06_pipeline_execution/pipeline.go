package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

// Stage - функция, которая принимает канал со значениями,
// выполняет над ними какую-то работу и возвращает канал с результатами
// этой функции.
type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Значения из канала in передаются в первый этап Stage,
	// результат которой - входной канал для следующего этапа Stage

	inNext := in // канал данных для следующего этапа

	for _, stage := range stages {
		inNext = stage(inNext)
	}

	// resChannel - канал с результатами пайплайна,
	// который мы закрываем по окончании работы
	// или при получении сигнала остановки.
	resChannel := make(chan interface{})

	go func() {
		defer close(resChannel)
		for {
			select {
			case v, ok := <-inNext:
				if !ok {
					return
				}
				resChannel <- v
			case <-done:
				return
			}
		}
	}()

	return resChannel
}
