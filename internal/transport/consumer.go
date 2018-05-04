package transport

import (
	"github.com/streadway/amqp"
	"taskd/internal/taskd"
)

type consumer struct {
	outStream     chan<- taskd.Request
	connector     dialer
	converter     Converter
	mainQueue     mainQueue
	consumeParams consumeParams
}

func newConsumer(outStream chan<- taskd.Request, connector dialer, converter Converter) *consumer {
	return &consumer{outStream: outStream, connector: connector, converter: converter}
}

func (consumer *consumer) run() {
	for {
		inputSteam, err := consumer.connect()
		if err != nil {
			break
		}
		consumer.handler(inputSteam)
	}
}

func (consumer *consumer) connect() (<-chan amqp.Delivery, error) {
	connection, err := consumer.connector.dial()
	if err != nil {
		return nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	if err := consumer.topology(channel); err != nil {
		return nil, err
	}
	inputStream, err := consumer.consume(channel)
	return inputStream, err
}

func (consumer *consumer) topology(channel *amqp.Channel) error {
	_, err := channel.QueueDeclare(
		consumer.mainQueue.Name,
		consumer.mainQueue.Durable,
		consumer.mainQueue.AutoDelete,
		consumer.mainQueue.Exclusive,
		consumer.mainQueue.NoWait,
		consumer.mainQueue.Args,
	)
	return err
}

func (consumer *consumer) consume(channel *amqp.Channel) (<-chan amqp.Delivery, error) {
	inputStream, err := channel.Consume(
		consumer.mainQueue.Name,
		consumer.consumeParams.Name,
		consumer.consumeParams.AutoAck,
		consumer.consumeParams.Exclusive,
		consumer.consumeParams.NoLocal,
		consumer.consumeParams.NoWait,
		consumer.consumeParams.Args,
	)
	return inputStream, err
}

func (consumer *consumer) handler(inputSteam <-chan amqp.Delivery) {
	for delivery := range inputSteam {
		request, err := consumer.converter.fromDelivery(&delivery)
		if err != nil {
			//packet drop
			continue
		}
		consumer.outStream <- *request
	}
}
