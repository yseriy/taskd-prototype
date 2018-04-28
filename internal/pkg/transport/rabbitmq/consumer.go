package rabbitmq

import (
	"github.com/streadway/amqp"
	"taskd/internal/pkg/taskd"
	"taskd/internal/pkg/transport/rabbitmq/connector"
)

type consumer struct {
	outStream chan<- taskd.Request
	connector connector.InputConnector
	converter converter
}

func newConsumer(outStream chan<- taskd.Request, connector connector.InputConnector, converter converter) *consumer {
	return &consumer{outStream: outStream, connector: connector, converter: converter}
}

func (consumer *consumer) run() {
	for {
		inputSteam, err := consumer.connector.Connect()
		if err != nil {
			break
		}
		consumer.handler(inputSteam)
	}
}

func (consumer *consumer) handler(inputSteam <-chan amqp.Delivery) {
	for delivery := range inputSteam {
		request := consumer.converter.fromDelivery(&delivery)
		consumer.outStream <- *request
	}
}
