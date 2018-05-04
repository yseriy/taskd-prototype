package converter

import (
	"github.com/streadway/amqp"
	"taskd/internal/taskd"
	"taskd/internal/transport"
)

type converter struct {
}

func (converter *converter) fromDelivery(delivery *amqp.Delivery) (*taskd.Request, error) {
	return nil, nil
}

func (converter *converter) toPublishing(response *taskd.Response) (transport.Address, *amqp.Publishing, error) {
	return nil, nil, nil
}

func New() transport.Converter {
	return &converter{}
}
