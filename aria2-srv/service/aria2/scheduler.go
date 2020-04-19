package aria2

import "github.com/zyxar/argo/rpc"

type Scheduler interface {
	WorkerReady(chan []rpc.Event)
	Submit([]rpc.Event)
	WorkerChannel() chan []rpc.Event
	Run()
}

// 异步调度器
type AsyncScheduler struct {
	eventChan  chan []rpc.Event
	workerChan chan chan []rpc.Event
}

func (s *AsyncScheduler) WorkerReady(c chan []rpc.Event) {
	s.workerChan <- c
}

func (s *AsyncScheduler) Submit(events []rpc.Event) {
	s.eventChan <- events
}

func (s *AsyncScheduler) WorkerChannel() chan []rpc.Event {
	return make(chan []rpc.Event)
}

func (s *AsyncScheduler) Run() {
	s.eventChan = make(chan []rpc.Event)
	s.workerChan = make(chan chan []rpc.Event)

	go func() {
		var eventQueue [][]rpc.Event
		var workerQueue []chan []rpc.Event

		for {
			var activeEvent []rpc.Event
			var activeWorker chan []rpc.Event

			if len(eventQueue) > 0 && len(workerQueue) > 0 {
				activeEvent = eventQueue[0]
				activeWorker = workerQueue[0]

				select {
				case events := <-s.eventChan:
					eventQueue = append(eventQueue, events)
				case worker := <-s.workerChan:
					workerQueue = append(workerQueue, worker)
				case activeWorker <- activeEvent:
					eventQueue = eventQueue[1:]
					workerQueue = workerQueue[1:]
				}
			}
		}
	}()
}
