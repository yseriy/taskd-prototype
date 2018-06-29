package rabbitmq

import "github.com/streadway/amqp"

type adapter struct {
	connection *amqp.Connection
}

func (a adapter) Channel() (amqpChannel, error) {
	return a.connection.Channel()
}

func (a adapter) Close() error {
	return a.connection.Close()
}
