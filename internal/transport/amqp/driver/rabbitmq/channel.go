package rabbitmq

import (
	"github.com/streadway/amqp"
	"taskd/internal/transport"
	"taskd/internal/transport/amqp/driver"
)

type channel struct {
	rabbitmqChannel rabbitmqChannel
}

func (c channel) QueueDeclare(queue *driver.AMQPQueue) error {
	_, err := c.rabbitmqChannel.QueueDeclare(queue.Name, queue.Durable, queue.AutoDelete, queue.Exclusive, false,
		amqp.Table(queue.Args))
	return mapError(err)
}

func (c channel) Publish(message *transport.Message) error {
	err := c.rabbitmqChannel.Publish("", message.Address, true, true, messageToPublishing(message))
	return mapError(err)
}

func messageToPublishing(message *transport.Message) amqp.Publishing {
	return amqp.Publishing{
		ContentType: message.ContentType,
		Body:        message.Body,
	}
}

func (c channel) Consume(queue *driver.AMQPQueue) (<-chan transport.Message, error) {
	deliveryStream, err := c.rabbitmqChannel.Consume(queue.Name, "", true, false, false,
		false, nil)
	if err != nil {
		return nil, mapError(err)
	}
	messageStream := make(chan transport.Message)
	go convertStream(deliveryStream, messageStream)
	return messageStream, nil
}

func convertStream(deliveryStream <-chan amqp.Delivery, messageStream chan transport.Message) {
	for {
		delivery, ok := <-deliveryStream
		if !ok {
			break
		}
		message := deliveryToMessage(delivery)
		messageStream <- message
	}
}

func deliveryToMessage(delivery amqp.Delivery) transport.Message {
	return transport.Message{
		ContentType: delivery.ContentType,
		Address:     delivery.ReplyTo,
		Body:        delivery.Body,
	}
}

func (c channel) Close() error {
	return c.rabbitmqChannel.Close()
}
