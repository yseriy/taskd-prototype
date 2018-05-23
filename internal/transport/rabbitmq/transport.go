package rabbitmq

import (
	"errors"
	"github.com/streadway/amqp"
	"taskd/internal/taskd"
)

type connector interface {
	Connect() channel
}

type channel interface {
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

type converter interface {
	FromDelivery(*amqp.Delivery) (*taskd.Request, error)
	ToPublishing(*taskd.Response) (string, string, *amqp.Publishing, error)
}

type rabbitmqTransport struct {
	config RabbitConfig
}

func New(config RabbitConfig) *rabbitmqTransport {
	return &rabbitmqTransport{config: config}
}

func (transport *rabbitmqTransport) Connect(requests chan<- taskd.Request, responses <-chan taskd.Response) {
	consumer := newConsumer(
		requests,
		newConnector(transport.config.Url, transport.config.Timeout),
		newConverter(),
		transport.config.MainQueue,
		transport.config.ConsumeParams,
	)
	publisher := newPublisher(
		responses,
		newConnector(transport.config.Url, transport.config.Timeout),
		newConverter(),
		transport.config.DestinationQueueTemplate,
		transport.config.PublishingParams,
	)

	go consumer.run()
	go publisher.run()
}

type consumer struct {
	outStream     chan<- taskd.Request
	connector     connector
	converter     converter
	mainQueue     MainQueue
	consumeParams ConsumeParams
}

func newConsumer(outStream chan<- taskd.Request, connector connector, converter converter,
	mainQueue MainQueue, consumeParams ConsumeParams) *consumer {
	return &consumer{
		outStream:     outStream,
		connector:     connector,
		converter:     converter,
		mainQueue:     mainQueue,
		consumeParams: consumeParams,
	}
}

func (consumer *consumer) run() {
	for {
		consumer.handler(consumer.connect())
	}
}

func (consumer *consumer) connect() <-chan amqp.Delivery {
	channel := consumer.connector.Connect()

	if err := consumer.topology(channel); err != nil {
		panic(err)
	}
	inputStream, err := consumer.consume(channel)
	if err != nil {
		panic(err)
	}
	return inputStream
}

func (consumer *consumer) topology(channel channel) error {
	_, err := channel.QueueDeclare(
		consumer.mainQueue.Name,
		consumer.mainQueue.Durable,
		consumer.mainQueue.AutoDelete,
		consumer.mainQueue.Exclusive,
		consumer.mainQueue.NoWait,
		consumer.mainQueue.Args,
	)
	return err
}

func (consumer *consumer) consume(channel channel) (<-chan amqp.Delivery, error) {
	inputStream, err := channel.Consume(
		consumer.mainQueue.Name,
		consumer.consumeParams.Name,
		consumer.consumeParams.AutoAck,
		consumer.consumeParams.Exclusive,
		consumer.consumeParams.NoLocal,
		consumer.consumeParams.NoWait,
		consumer.consumeParams.Args,
	)
	return inputStream, err
}

func (consumer *consumer) handler(inputSteam <-chan amqp.Delivery) {
	for delivery := range inputSteam {
		request, err := consumer.converter.FromDelivery(&delivery)
		if err != nil {
			//packet drop
			continue
		}
		consumer.outStream <- *request
	}
}

type publisher struct {
	inStream         <-chan taskd.Response
	connector        connector
	converter        converter
	queueTemplate    DestinationQueueTemplate
	publishingParams PublishingParams
}

func newPublisher(inStream <-chan taskd.Response, connector connector, converter converter,
	queueTemplate DestinationQueueTemplate, publishingParams PublishingParams) *publisher {
	return &publisher{
		inStream:         inStream,
		connector:        connector,
		converter:        converter,
		queueTemplate:    queueTemplate,
		publishingParams: publishingParams,
	}
}

func (publisher *publisher) run() {
	for {
		publisher.handler(publisher.connect())
	}
}

func (publisher *publisher) connect() (*publisherHelper) {
	channel := publisher.connector.Connect()
	return newPublisherHelper(channel, publisher.queueTemplate, publisher.publishingParams)
}

func (publisher *publisher) handler(helper *publisherHelper) {
	for {
		response := publisher.receive()

		exchange, routingKey, publishing, err := publisher.converter.ToPublishing(response)
		if err != nil {
			//packet drop
			continue
		}
		if err := helper.send(exchange, routingKey, publishing); err != nil {
			break
		}
	}
}

func (publisher *publisher) receive() *taskd.Response {
	response, ok := <-publisher.inStream
	if !ok {
		panic(errors.New("response channel closed"))
	}
	return &response
}

type publisherHelper struct {
	channel          channel
	queueTemplate    DestinationQueueTemplate
	publishingParams PublishingParams
}

func newPublisherHelper(channel channel, queueTemplate DestinationQueueTemplate,
	publishingParams PublishingParams) *publisherHelper {
	return &publisherHelper{channel: channel, queueTemplate: queueTemplate, publishingParams: publishingParams}
}

func (helper *publisherHelper) send(exchange, routingKey string, publishing *amqp.Publishing) error {
	if err := helper.topology(routingKey); err != nil {
		return err
	}
	return helper.publish(exchange, routingKey, publishing)
}

func (helper *publisherHelper) topology(queueName string) error {
	_, err := helper.channel.QueueDeclare(
		queueName,
		helper.queueTemplate.Durable,
		helper.queueTemplate.AutoDelete,
		helper.queueTemplate.Exclusive,
		helper.queueTemplate.NoWait,
		helper.queueTemplate.Args,
	)
	return err
}

func (helper *publisherHelper) publish(exchange, routingKey string, publishing *amqp.Publishing) error {
	err := helper.channel.Publish(
		exchange,
		routingKey,
		helper.publishingParams.Mandatory,
		helper.publishingParams.Immediate,
		*publishing,
	)
	return err
}
