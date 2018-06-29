package rabbitmq

import (
	"github.com/streadway/amqp"
	"taskd/internal/transport"
)

type (
	connection interface {
		Channel() (channel, error)
	}

	channel interface {
		QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
		Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
		Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool,
			args amqp.Table) (<-chan amqp.Delivery, error)
	}

	closer interface {
		Close() error
	}
)

type (
	amqpConnection interface {
		Channel() (amqpChannel, error)
		closer
	}

	amqpChannel interface {
		channel
		closer
	}
)

type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	Args       map[string]interface{}
}

func mapConnector(connector func(url string) (*amqp.Connection, error)) func(url string) (amqpConnection, error) {
	return func(url string) (amqpConnection, error) {
		connection, err := connector(url)
		if err != nil {
			return nil, err
		}
		return &adapter{connection: connection}, nil
	}
}

func toMessage(delivery *amqp.Delivery) *transport.Message {
	return &transport.Message{
		Address:     delivery.ReplyTo,
		ContentType: delivery.ContentType,
		Body:        delivery.Body,
	}
}

func declareParams(queue Queue) (string, bool, bool, bool, bool, amqp.Table) {
	return queue.Name, queue.Durable, queue.AutoDelete, queue.Exclusive, false, amqp.Table(queue.Args)
}

func publishParams(message transport.Message) (string, string, bool, bool, amqp.Publishing) {
	return "", message.Address, false, true, amqp.Publishing{
		ContentType: message.ContentType,
		Body:        message.Body,
	}
}

func consumParams(queue Queue) (string, string, bool, bool, bool, bool, amqp.Table) {
	return queue.Name, "", true, false, false, false, amqp.Table(queue.Args)
}

const (
	accessRefused      = 403
	preconditionFailed = 406
)

func isDisconnect(err error) bool {
	amqpError := err.(amqp.Error)
	switch amqpError.Code {
	case accessRefused:
		return false
	case preconditionFailed:
		return false
	default:
		return true
	}
}

func isNotDisconnect(err error) bool {
	return !isDisconnect(err)
}
