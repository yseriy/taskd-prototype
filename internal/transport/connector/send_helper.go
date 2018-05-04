package connector

import (
	"github.com/streadway/amqp"
	"taskd/internal/transport"
)

type sendHelper struct {
	channel          *amqp.Channel
	queueTemplate    destinationQueueTemplate
	publishingParams publishingParams
}

func newSendHelper(channel *amqp.Channel, queueTemplate destinationQueueTemplate,
	publishingParams publishingParams) transport.SendHelper {
	return &sendHelper{channel: channel, queueTemplate: queueTemplate, publishingParams: publishingParams}
}

func (publisher *sendHelper) Send(address transport.Address, publishing *amqp.Publishing) error {
	if err := publisher.topology(address.RoutingKey); err != nil {
		return err
	}
	err := publisher.publish(address.Exchange, address.RoutingKey, publishing)
	return err
}

func (publisher *sendHelper) topology(queueName string) error {
	_, err := publisher.channel.QueueDeclare(
		queueName,
		publisher.queueTemplate.Durable,
		publisher.queueTemplate.AutoDelete,
		publisher.queueTemplate.Exclusive,
		publisher.queueTemplate.NoWait,
		publisher.queueTemplate.Args,
	)
	return err
}

func (publisher *sendHelper) publish(exchange, routingKey string, publishing *amqp.Publishing) error {
	err := publisher.channel.Publish(
		exchange,
		routingKey,
		publisher.publishingParams.Mandatory,
		publisher.publishingParams.Immediate,
		*publishing,
	)
	return err
}
