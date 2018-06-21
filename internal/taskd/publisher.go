package taskd

import "taskd/internal/taskd"

type publisher struct {
	queue        amqpQueue
	connector    connector
	converter    converter
	inputStream  chan amqpRequest
	outputStream chan taskd.Request
}

func (p *publisher) Connect() (chan taskd.Request, error) {
	inputStream, err := p.setupInputStream()
	if err != nil {
		return nil, err
	}
	p.inputStream = inputStream
	p.outputStream = make(chan taskd.Request)
	return p.outputStream, nil
}

func (p *publisher) handler() {
	for {
		amqpRequest, err := p.receive()
		if err != nil {
			break
		}
		request, err := p.converter.fromTransport(amqpRequest)
		if err != nil {
			continue // drop request
		}
		p.outputStream <- request
	}
}

func (p *publisher) receive() (*amqpRequest, error) {
	for {
		request, ok := <-p.inputStream
		if ok {
			return &request, nil
		}
		inputStream, err := p.setupInputStream()
		if err != nil {
			return nil, err
		}
		p.inputStream = inputStream
	}
}

func (p publisher) setupInputStream() (chan amqpRequest, error) {
	for {
		channel, err := p.connector.Connect()
		if err != nil {
			return nil, err
		}
		if err := channel.QueueDeclare(&p.queue); err != nil {
			if channel.IsDisconnect(err) {
				continue
			} else {
				return nil, err
			}
		}
		if inputStream, err := channel.Consume(); err == nil {
			return inputStream, nil
		}
	}
}
