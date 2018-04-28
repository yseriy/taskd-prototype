package connector

import "github.com/streadway/amqp"

type PublisherHelper struct {
	channel          *amqp.Channel
	queueTemplate    DestinationQueueTemplate
	publishingParams PublishingParams
}

func newPublisherHelper(channel *amqp.Channel, queueTemplate DestinationQueueTemplate,
	publishingParams PublishingParams) *PublisherHelper {
	return &PublisherHelper{channel: channel, queueTemplate: queueTemplate, publishingParams: publishingParams}
}

func (publisher *PublisherHelper) Send(stash *amqpTransportStash, publishing *amqp.Publishing) error {
	if err := publisher.topology(stash); err != nil {
		return err
	}
	err := publisher.publish(stash, publishing)
	return err
}

func (publisher *PublisherHelper) topology(stash *amqpTransportStash) error {
	_, err := publisher.channel.QueueDeclare(
		stash.destinationQueueName,
		publisher.queueTemplate.Durable,
		publisher.queueTemplate.AutoDelete,
		publisher.queueTemplate.Exclusive,
		publisher.queueTemplate.NoWait,
		publisher.queueTemplate.Args,
	)
	return err
}

func (publisher *PublisherHelper) publish(stash *amqpTransportStash, publishing *amqp.Publishing) error {
	err := publisher.channel.Publish(
		stash.destinationAddress.exchange,
		stash.destinationAddress.routingKey,
		publisher.publishingParams.Mandatory,
		publisher.publishingParams.Immediate,
		*publishing,
	)
	return err
}
