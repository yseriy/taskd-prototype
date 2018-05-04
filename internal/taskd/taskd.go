package taskd

import (
	"sync"
	"taskd/internal/transport"
)

type taskd struct {
	requestStream  chan Request
	responseStream chan Response
	transport      Transport
}

func NewTaskd() *taskd {
	return &taskd{
		requestStream:  make(chan Request),
		responseStream: make(chan Response),
		transport: transport.New()
	}
}

func (taskd *taskd) Run() {
	taskd.transport.Connect(taskd.requestStream, taskd.responseStream)
	lock()
}

func lock() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
