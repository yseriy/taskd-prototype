package transport

import (
	"github.com/streadway/amqp"
	"taskd/internal/taskd"
)

type consumer struct {
	outStream chan<- taskd.Request
	connector InputConnector
	converter Converter
}

func newConsumer(outStream chan<- taskd.Request, connector InputConnector, converter Converter) *consumer {
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
		request, err := consumer.converter.fromDelivery(&delivery)
		if err != nil {
			//packet drop
			continue
		}
		consumer.outStream <- *request
	}
}
