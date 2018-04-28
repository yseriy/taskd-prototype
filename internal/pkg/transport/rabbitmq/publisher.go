package rabbitmq

import (
	"errors"
	"taskd/internal/pkg/taskd"
	"taskd/internal/pkg/transport/rabbitmq/connector"
)

type publisher struct {
	inStream  <-chan taskd.Response
	connector connector.OutputConnector
	converter converter
}

func newPublisher(inStream <-chan taskd.Response, connector connector.OutputConnector, converter converter) *publisher {
	return &publisher{inStream: inStream, connector: connector, converter: converter}
}

func (publisher *publisher) run() {
	for {
		helper := publisher.connector.Connect()
		publisher.handler(helper)
	}
}

func (publisher *publisher) handler(helper *connector.PublisherHelper) {
	for {
		response := publisher.receive()
		stash, publishing := publisher.converter.fromResponse(response)
		if err := helper.Send(stash, publishing); err != nil {
			break
		}
	}
}

func (publisher *publisher) receive() *taskd.Response {
	response, ok := <-publisher.inStream
	if !ok {
		panic(errors.New("response channel closed"))
	}
	return &response
}
