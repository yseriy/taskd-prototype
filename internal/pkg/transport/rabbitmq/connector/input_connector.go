package connector

import "github.com/streadway/amqp"

type InputConnector struct {
	connector     connector
	mainQueue     MainQueue
	consumeParams ConsumeParams
}

func newInputConnector(connector connector, mainQueue MainQueue, consumeParams ConsumeParams) *InputConnector {
	return &InputConnector{connector: connector, mainQueue: mainQueue, consumeParams: consumeParams}
}

func (inputConnector *InputConnector) Connect() (<-chan amqp.Delivery, error) {
	channel := inputConnector.connector.connect()
	if err := inputConnector.topology(channel); err != nil {
		return nil, err
	}
	inputStream, err := inputConnector.consume(channel)
	return inputStream, err
}

func (inputConnector *InputConnector) topology(channel *amqp.Channel) error {
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

func (inputConnector *InputConnector) consume(channel *amqp.Channel) (<-chan amqp.Delivery, error) {
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
