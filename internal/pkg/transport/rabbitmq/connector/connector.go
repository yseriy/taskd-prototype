package connector

import (
	"github.com/streadway/amqp"
	"time"
)

type connector struct {
	url        string
	timeout    int
	connection *amqp.Connection
}

func newConnector(url string, timeout int) *connector {
	return &connector{url: url, timeout: timeout}
}

func (connector *connector) connect() *amqp.Channel {
	connector.reset()
	for {
		connection, channel, err := connector.dial()
		if err != nil {
			time.Sleep(time.Duration(connector.timeout) * time.Second)
			continue
		}
		connector.connection = connection
		return channel
	}
}

func (connector *connector) reset() {
	if connector.connection != nil {
		connector.connection.Close()
		connector.connection = nil
	}
}

func (connector *connector) dial() (*amqp.Connection, *amqp.Channel, error) {
	connection, err := amqp.Dial(connector.url)
	if err != nil {
		return connection, nil, err
	}
	channel, err := connection.Channel()
	return connection, channel, err
}
