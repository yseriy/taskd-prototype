package rabbitmq

import "github.com/streadway/amqp"

type mockRabbitMQConnection struct {
	channel func() (rabbitmqChannel, error)
	close   func() error
}

func (c *mockRabbitMQConnection) Channel() (rabbitmqChannel, error) {
	if c.channel != nil {
		return c.channel()
	}
	return nil, nil
}

func (c *mockRabbitMQConnection) Close() error {
	if c.close != nil {
		return c.close()
	}
	return nil
}

type mockRabbitMQChannel struct {
}

func (c *mockRabbitMQChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

func (c *mockRabbitMQChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return nil, nil
}

func (c *mockRabbitMQChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return nil
}

func (c *mockRabbitMQChannel) Close() error {
	return nil
}
