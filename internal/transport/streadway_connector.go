package transport

import "github.com/streadway/amqp"

type streadwayConnector struct {
	url string
}

func (c streadwayConnector) Connect() (connection, error) {
	connection, err := amqp.Dial(c.url)
	return &streadwayConnection{connection: connection}, err
}

type streadwayConnection struct {
	connection *amqp.Connection
}

func (c streadwayConnection) Channel() (channel, error) {
	channel, err := c.connection.Channel()
	return &streadwayChannel{channel: channel}, err
}

func (c streadwayConnection) Close() error {
	return c.connection.Close()
}

type streadwayChannel struct {
	channel *amqp.Channel
}

func (c streadwayChannel) QueueDeclare(queue amqpQueue) error {
	_, err := c.channel.QueueDeclare(
		queue.Name,
		queue.Durable,
		queue.AutoDelete,
		queue.Exclusive,
		queue.NoWait,
		declareMap(queue.Args),
	)
	return err
}

func declareMap(table table) amqp.Table {
	return amqp.Table(table)
}

func (c streadwayChannel) Publish(response amqpResponse) error {
	return c.channel.Publish(
		response.Exchange,
		response.RoutingKey,
		response.Mandatory,
		response.Immediate,
		responseMap(response),
	)
}

func responseMap(response amqpResponse) amqp.Publishing {
	return amqp.Publishing{Body: response.Body}
}
