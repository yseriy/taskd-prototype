package rabbitmq

import (
	"github.com/streadway/amqp"
	"time"
)

type amqpConnector struct {
}

func (connector *amqpConnector) Dial() connection {
	return &amqpConnection{}
}

type amqpConnection struct {
}

func (connection *amqpConnection) Channel() (channel, error) {
	return &amqpChannel{}, nil
}

type amqpChannel struct {
}

func (channel *amqpChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

func (channel *amqpChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool,
	args amqp.Table) (<-chan amqp.Delivery, error) {
	return nil, nil
}

func (channel *amqpChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return nil
}

type repeater struct {
	url     string
	timeout int
}

func newConnector(url string, timeout int) *repeater {
	return &repeater{url: url, timeout: timeout}
}

func (repeater *repeater) Dial() *amqp.Connection {
	for {
		connection, err := amqp.Dial(repeater.url)
		if err != nil {
			time.Sleep(time.Duration(repeater.timeout) * time.Second)
			continue
		}
		return connection
	}
}
