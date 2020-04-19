package aria2

import "github.com/zyxar/argo/rpc"

type Engine interface {
	Init(int)
	Submit([]rpc.Event)
}

type AsyncEngine struct {
	s Scheduler
}

func NewAsyncEngine() *AsyncEngine {
	return &AsyncEngine{s: new(AsyncScheduler)}
}

func (a *AsyncEngine) Submit(events []rpc.Event) {
	a.s.Submit(events)
}

func (a *AsyncEngine) Init(c int) {
	a.s.Run()
	for i := 0; i < c; i++ {
		createWorker(a.s.WorkerChannel(), a.s)
	}
}
