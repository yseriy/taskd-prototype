package rabbitmq

import (
	"github.com/streadway/amqp"
	"taskd/internal/transport/amqp/driver"
)

const (
	accessRefused      = 403
	preconditionFailed = 406
)

type connector struct {
	connect func(url string) (*amqp.Connection, error)
	url     string
}

func NewConnector(url string) *connector {
	return &connector{connect: amqp.Dial, url: url}
}

func (a connector) Dial() (driver.AMQPConnection, error) {
	rabbitmqConnection, err := mapConnection(a.connect, a.url)
	if err != nil {
		return nil, mapError(err)
	}
	return &connection{rabbitmqConnection: rabbitmqConnection}, nil
}

func mapError(err error) error {
	amqpError := err.(amqp.Error)
	switch amqpError.Code {
	case accessRefused:
		return amqpError
	case preconditionFailed:
		return amqpError
	default:
		return driver.ErrorDisconnect{DriverError: amqpError}
	}
}
