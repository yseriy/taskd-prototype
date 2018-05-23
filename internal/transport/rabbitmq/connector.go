package rabbitmq

import (
	"github.com/streadway/amqp"
	"time"
)

type amqpConnector struct {
	url        string
	connection *amqp.Connection
}

func (connector *amqpConnector) Connect() (*amqp.Channel, error) {
	connector.Close()
	connection, err := amqp.Dial(connector.url)
	if err != nil {
		return nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	connector.connection = connection
	return channel, nil
}

func (connector *amqpConnector) Close() {
	if connector.connection != nil {
		connector.connection.Close()
	}
}

type amqpRepeater struct {
	connector amqpConnector
	timeout   int
}

func newConnector(url string, timeout int) *amqpRepeater {
	return &amqpRepeater{connector: amqpConnector{url: url}, timeout: timeout}
}

func (repeater *amqpRepeater) Connect() channel {
	for {
		channel, err := repeater.connector.Connect()
		if err != nil {
			time.Sleep(time.Duration(repeater.timeout) * time.Second)
			continue
		}
		return channel
	}
}
