package rabbitmq

import (
	"github.com/streadway/amqp"
)

type (
	rabbitmqConnection interface {
		Channel() (rabbitmqChannel, error)
		Close() error
	}

	rabbitmqChannel interface {
		QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
		Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
		Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
		Close() error
	}
)

type adapter struct {
	connection *amqp.Connection
}

func (a adapter) Channel() (rabbitmqChannel, error) {
	return a.connection.Channel()
}

func (a adapter) Close() error {
	return a.connection.Close()
}

func mapConnection(connect func(url string) (*amqp.Connection, error), url string) (rabbitmqConnection, error) {
	connection, err := connect(url)
	if err != nil {
		return nil, err
	}
	return adapter{connection: connection}, nil
}
