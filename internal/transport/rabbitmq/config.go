package rabbitmq

import "github.com/streadway/amqp"

type RabbitConfig struct {
	Url                      string
	Timeout                  int
	MainQueue                MainQueue
	ConsumeParams            ConsumeParams
	DestinationQueueTemplate DestinationQueueTemplate
	PublishingParams         PublishingParams
}

type MainQueue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type ConsumeParams struct {
	Name      string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

type DestinationQueueTemplate struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type PublishingParams struct {
	Mandatory, Immediate bool
}
