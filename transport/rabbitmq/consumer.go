package rabbitmq

import (
	"github.com/streadway/amqp"
	"taskd"
)

type Consumer struct {
	in  <-chan amqp.Delivery
	out chan<- taskd.Request
}

func NewConsumer(in <-chan amqp.Delivery, out chan<- taskd.Request) *Consumer {
	return &Consumer{in: in, out: out}
}

func (consumer *Consumer) Run() {
	for {
		if receiveStatus := consumer.handler(); receiveStatus == stop {
			break
		}
	}
}

func (consumer *Consumer) handler() receiveStatus {
	delivery, ok := consumer.receive()
	if !ok {
		return stop
	}
	request, err := consumer.convert(delivery)
	if err != nil {
		return resume
	}
	consumer.send(request)
	return resume
}

func (consumer *Consumer) receive() (*amqp.Delivery, bool) {
	delivery, ok := <-consumer.in
	return &delivery, ok
}

func (consumer *Consumer) convert(delivery *amqp.Delivery) (*taskd.Request, error) {
	body := string(delivery.Body)
	return &taskd.Request{Body: body}, nil
}

func (consumer *Consumer) send(request *taskd.Request) {
	consumer.out <- *request
}
