package transport

import (
	"errors"
	"taskd/internal/taskd"
)

type publisher struct {
	inStream         <-chan taskd.Response
	connector        dialer
	converter        Converter
	queueTemplate    destinationQueueTemplate
	publishingParams publishingParams
}

func newPublisher(inStream <-chan taskd.Response, connector dialer, converter Converter) *publisher {
	return &publisher{inStream: inStream, connector: connector, converter: converter}
}

func (publisher *publisher) run() {
	for {
		helper, err := publisher.connect()
		if err != nil {
			break
		}
		publisher.handler(helper)
	}
}

func (publisher *publisher) connect() (*sendHelper, error) {
	connection, err := publisher.connector.dial()
	if err != nil {
		return nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	return newSendHelper(channel, publisher.queueTemplate, publisher.publishingParams), nil
}

func (publisher *publisher) handler(helper SendHelper) {
	for {
		response := publisher.receive()
		address, publishing, err := publisher.converter.toPublishing(response)
		if err != nil {
			//packet drop
			continue
		}
		if err := helper.Send(address, publishing); err != nil {
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
