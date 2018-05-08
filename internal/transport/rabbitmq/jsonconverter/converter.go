package jsonconverter

import (
	"github.com/streadway/amqp"
	"taskd/internal/taskd"
)

func New() *JsonConverter {
	return &JsonConverter{}
}

type JsonConverter struct {
}

func (converter *JsonConverter) FromDelivery(delivery *amqp.Delivery) (*taskd.Request, error) {
	return nil, nil
}

func (converter *JsonConverter) ToPublishing(response *taskd.Response) (string, string, *amqp.Publishing, error) {
	return "", "", &amqp.Publishing{}, nil
}
