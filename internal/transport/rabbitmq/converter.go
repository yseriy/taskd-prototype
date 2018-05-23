package rabbitmq

import (
	"github.com/streadway/amqp"
	"taskd/internal/taskd"
)

func newConverter() *jsonConverter {
	return &jsonConverter{}
}

type jsonConverter struct {
}

func (converter *jsonConverter) FromDelivery(delivery *amqp.Delivery) (*taskd.Request, error) {
	return &taskd.Request{}, nil
}

func (converter *jsonConverter) ToPublishing(response *taskd.Response) (string, string, *amqp.Publishing, error) {
	return "", "", &amqp.Publishing{}, nil
}
