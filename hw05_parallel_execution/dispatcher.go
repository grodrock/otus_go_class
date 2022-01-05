package hw05parallelexecution

type Dispatcher struct {
	workersPool   []*Worker
	TaskChan      chan Task
	QueueWorkers  chan chan Task
	ErrorsChannel chan error
	End           chan bool
	ErrorsCount   int
}

func (d *Dispatcher) StartWorkers() {
	l := len(d.workersPool)
	for i := 0; i < l; i++ {
		w := Worker{
			id:             i,
			WorkersChannel: d.QueueWorkers,
			WorkChannel:    make(chan Task),
			errorChannel:   d.ErrorsChannel,
			end:            make(chan bool),
		}
		w.Start()
		d.workersPool[i] = &w
	}
	go d.process()
}

func (d *Dispatcher) process() {
	for {
		select {
		case task := <-d.TaskChan:
			workerC := <-d.QueueWorkers
			workerC <- task
		case res := <-d.ErrorsChannel:
			if res != nil {
				d.ErrorsCount++
			}
		case <-d.End:
			return
		}
	}
}

func (d *Dispatcher) StopWorkers() {
	for _, w := range d.workersPool {
		w.Stop()
	}
}

func (d *Dispatcher) AddTaskToProcess(task Task) {
	d.TaskChan <- task
}
