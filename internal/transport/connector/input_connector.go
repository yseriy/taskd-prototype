package connector

import (
	"github.com/streadway/amqp"
	"taskd/internal/transport"
)

type inputConnector struct {
	connector     dialer
	mainQueue     mainQueue
	consumeParams consumeParams
}

func newInputConnector(connector dialer, mainQueue mainQueue, consumeParams consumeParams) transport.InputConnector {
	return &inputConnector{connector: connector, mainQueue: mainQueue, consumeParams: consumeParams}
}

func (inputConnector *inputConnector) Connect() (<-chan amqp.Delivery, error) {
	connection, err := inputConnector.connector.dial()
	if err != nil {
		return nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	if err := inputConnector.topology(channel); err != nil {
		return nil, err
	}
	inputStream, err := inputConnector.consume(channel)
	return inputStream, err
}

func (inputConnector *inputConnector) topology(channel *amqp.Channel) error {
	_, err := channel.QueueDeclare(
		inputConnector.mainQueue.Name,
		inputConnector.mainQueue.Durable,
		inputConnector.mainQueue.AutoDelete,
		inputConnector.mainQueue.Exclusive,
		inputConnector.mainQueue.NoWait,
		inputConnector.mainQueue.Args,
	)
	return err
}

func (inputConnector *inputConnector) consume(channel *amqp.Channel) (<-chan amqp.Delivery, error) {
	inputStream, err := channel.Consume(
		inputConnector.mainQueue.Name,
		inputConnector.consumeParams.Name,
		inputConnector.consumeParams.AutoAck,
		inputConnector.consumeParams.Exclusive,
		inputConnector.consumeParams.NoLocal,
		inputConnector.consumeParams.NoWait,
		inputConnector.consumeParams.Args,
	)
	return inputStream, err
}
