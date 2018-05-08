package taskd

import (
	"sync"
)

type taskd struct {
	requestStream  chan Request
	responseStream chan Response
	transport      transport
}

func NewTaskd() *taskd {
	return &taskd{
		requestStream:  make(chan Request),
		responseStream: make(chan Response),
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
