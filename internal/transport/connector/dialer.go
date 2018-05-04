package connector

import (
	"github.com/streadway/amqp"
)

type dialer interface {
	dial() (*amqp.Connection, error)
}

type dialerImpl struct {
	url string
}

func newDialer(url string) dialer {
	return &dialerImpl{url: url}
}

func (dialer *dialerImpl) dial() (*amqp.Connection, error) {
	connection, err := amqp.Dial(dialer.url)
	return connection, err
}
