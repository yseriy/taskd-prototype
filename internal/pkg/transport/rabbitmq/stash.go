package rabbitmq

type amqpTransportStash struct {
	destinationQueueName string
	destinationAddress   amqpAddress
}

type amqpAddress struct {
	exchange, routingKey string
}
