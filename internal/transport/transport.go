package transport

import (
	"taskd/internal/taskd"
	"taskd/internal/transport/connector"
	"taskd/internal/transport/converter"
)

type transport struct {
	config *taskd.RabbitConfig
}

func (transport *transport) Connect(inStream chan<- taskd.Request, outStream <-chan taskd.Response) {
	cnt := connector.New(transport.config)
	cnv := converter.New()
	consumer := newConsumer(inStream, cnt.Input(), cnv)
	publisher := newPublisher(outStream, cnt.Output(), cnv)

	go consumer.run()
	go publisher.run()
}

func New(config *taskd.RabbitConfig) taskd.Transport {
	return &transport{config: config}
}
