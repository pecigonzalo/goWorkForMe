package Gworker

// Dispatcher is a work pool handler
type Dispatcher struct {
	WorkerPool  chan chan Job // A pool of workers channels that are registered with the dispatcher
	JobQueue    chan Job      // JobQueue is a buffered channel that we can send work requests on
	PriJobQueue chan Job      // PrioJobQueue is a buffered channel that we can send work requests on
}

// NewDispatcher Create a new dispatcher
func NewDispatcher(
	maxWorkers int,
	maxQueue int,
	maxPrioQueue int,
) *Dispatcher {
	return &Dispatcher{
		WorkerPool:  make(chan chan Job, maxWorkers),
		JobQueue:    make(chan Job, maxQueue),
		PriJobQueue: make(chan Job, maxPrioQueue),
	}
}

// Auto generate workers and Run
func (d *Dispatcher) Auto() {
	// starting n number of workers
	for i := 0; i < cap(d.WorkerPool); i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.Run()
}

// Run Worker pool handler
func (d *Dispatcher) Run() {
	sendToWorker := func(job Job) {
		// a job request has been received
		// try to obtain a worker job channel that is available.
		// this will block until a worker is idle
		jobChannel := <-d.WorkerPool

		// dispatch the job to the worker job channel
		jobChannel <- job
	}

	for {
		select {
		case job := <-d.PriJobQueue:
			sendToWorker(job)
		case job := <-d.JobQueue:
			sendToWorker(job)
		}
	}
}
