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

// getResultChannel - функция, которая возвращает канал с результатами пайплайна,
// закрывает его по окончании или при получении сигнала остановки.
func getResultChannel(in In, done In) Out {
	resChannel := make(chan interface{})

	go func() {
		for {
			select {
			case v, ok := <-in:
				if !ok {
					close(resChannel)
					return
				}
				resChannel <- v
			case <-done:
				close(resChannel)
				return
			}
		}
	}()

	return resChannel
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Значения из канала in передаются в первый этап Stage,
	// результат которой - входной канал для следующего этапа Stage

	inNext := in // канал данных для следующего этапа

	for _, stage := range stages {
		inNext = stage(inNext)
	}

	resultChannel := getResultChannel(inNext, done)

	return resultChannel
}
