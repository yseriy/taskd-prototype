package driver

import (
	"taskd/internal/transport"
	"github.com/streadway/amqp"
)

type (
	AMQPSender interface {
		AMQPDeclarer
		AMQPPublisher
	}

	AMQPDeclarer interface {
		QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	}

	AMQPPublisher interface {
		Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	}

	AMQPConsumer interface {
		Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool,
			args amqp.Table) (<-chan amqp.Delivery, error)
	}

	AMQPCloser interface {
		Close() error
	}
)

type (
	AMQPConnector interface {
		Dial() (AMQPConnection, error)
	}

	AMQPConnection interface {
		Channel() (AMQPChannel, error)
		Close() error
	}

	AMQPChannel interface {
		QueueDeclare(queue *AMQPQueue) error
		Publish(message *transport.Message) error
		Consume(queue *AMQPQueue) (<-chan transport.Message, error)
		Close() error
	}

	AMQPQueue struct {
		Name       string
		Durable    bool
		AutoDelete bool
		Exclusive  bool
		Args       map[string]interface{}
	}

	ErrorDisconnect struct {
		DriverError error
	}
)

func (e ErrorDisconnect) Error() string {
	return e.DriverError.Error()
}

func IsDisconnect(err error) bool {
	_, ok := err.(ErrorDisconnect)
	if !ok {
		return false
	}
	return true
}
