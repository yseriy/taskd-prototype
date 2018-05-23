package transport

import (
	"log"
)

type consumer struct {
	dialer      dialer
	inputStream chan interface{}
}

func (c consumer) Run() {
	for {
		connection := c.dialer.Dial()
		c.setupHandler(connection)
	}
}

func (c consumer) setupHandler(connection connection) {
	defer connection.Close()
	for {
		channel, err := connection.Channel()
		if err != nil {
			log.Print(err)
			break
		}
		c.handler(channel)
	}
}

func (c consumer) handler(channel channel) {
	for task := range c.inputStream {

	}
}

func (c consumer) send(channel channel, response amqpResponse) error {
	queue := amqpQueue{Name: response.RoutingKey}
	err := channel.QueueDeclare(queue)
	if err != nil {
		return err
	}
	return channel.Publish(response)
}
