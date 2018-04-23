package rabbitmq

import (
	"github.com/streadway/amqp"
	"taskd"
)

type Publisher struct {
	amqpChannel *amqp.Channel
	in          <-chan taskd.Response
	signal      <-chan struct{}
}

func NewPublisher(amqpChannel *amqp.Channel, in <-chan taskd.Response, signal <-chan struct{}) *Publisher {
	return &Publisher{amqpChannel: amqpChannel, in: in, signal: signal}
}

func (publisher *Publisher) Run() {
	for {
		if receiveStatus := publisher.receive(); receiveStatus == stop {
			break
		}
	}
}

func (publisher *Publisher) receive() receiveStatus {
	select {
	case response := <-publisher.in:
		return publisher.topology(response)
	case <-publisher.signal:
		return stop
	}
}

func (publisher *Publisher) topology(response taskd.Response) receiveStatus {
	_, err := publisher.amqpChannel.QueueDeclare("rr", true, false, false,
		false, nil)
	if err != nil {
		return resume
	}
	return publisher.send(response)
}

func (publisher *Publisher) send(response taskd.Response) receiveStatus {
	publishing := amqp.Publishing{ContentType: "text/plain", Body: []byte(response.Body)}
	err := publisher.amqpChannel.Publish("", "rr", false, false, publishing)
	if err != nil {
		return stop
	}
	return resume
}
