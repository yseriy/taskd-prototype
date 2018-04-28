package rabbitmq

import (
	"github.com/streadway/amqp"
	"taskd/internal/pkg/taskd"
)

type converter interface {
	fromDelivery(*amqp.Delivery) *taskd.Request
	fromResponse(*taskd.Response) (*amqpTransportStash, *amqp.Publishing)
}

type converterImpl struct {
}

func newConverter() converter {
	return &converterImpl{}
}

func (converter *converterImpl) fromDelivery(delivery *amqp.Delivery) *taskd.Request {
	return nil
}

func (converter *converterImpl) fromResponse(response *taskd.Response) (*amqpTransportStash, *amqp.Publishing) {
	stash := response.TransportStash.(amqpTransportStash)
	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(response.Body),
	}
	return &stash, &publishing
}
