package taskd

import "taskd/internal/taskd"

type converter interface {
	fromTask(taskd.Response) (*amqpResponse, error)
	fromTransport(*amqpRequest) (taskd.Request, error)
}
