package transport

import (
	"github.com/streadway/amqp"
	"taskd/internal/taskd"
)

type Converter interface {
	fromDelivery(*amqp.Delivery) (*taskd.Request, error)
	toPublishing(*taskd.Response) (Address, *amqp.Publishing, error)
}
