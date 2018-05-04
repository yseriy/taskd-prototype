package transport

import "github.com/streadway/amqp"

type Connector interface {
	Input() InputConnector
	Output() OutputConnector
}

type InputConnector interface {
	Connect() (<-chan amqp.Delivery, error)
}

type OutputConnector interface {
	Connect() (SendHelper, error)
}

type SendHelper interface {
	Send(Address, *amqp.Publishing) error
}
