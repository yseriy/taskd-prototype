package rabbitmq

import (
	"taskd/internal/transport/amqp/driver"
)

type connection struct {
	rabbitmqConnection rabbitmqConnection
	newChannel         func(rabbitmqChannel rabbitmqChannel) *channel
	mapError           func(err error) error
}

func (c connection) Channel() (driver.AMQPChannel, error) {
	rabbitmqChannel, err := c.rabbitmqConnection.Channel()
	if err != nil {
		return nil, c.mapError(err)
	}
	return c.newChannel(rabbitmqChannel), nil
}

func (c connection) Close() error {
	return c.rabbitmqConnection.Close()
}
