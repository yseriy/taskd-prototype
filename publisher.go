package taskd

import "github.com/streadway/amqp"

type Publisher struct {
	amqpChannel   *amqp.Channel
	outChannel    <-chan string
	signalChannel <-chan struct{}
}

func NewPublisher(amqpChannel *amqp.Channel, outChannel <-chan string, signalChannel <-chan struct{}) *Publisher {
	return &Publisher{amqpChannel: amqpChannel, outChannel: outChannel, signalChannel: signalChannel}
}

func (publisher *Publisher) Run() {
	go publisher.handel()
}

func (publisher *Publisher) handel() {
	for {
		if status := publisher.receive(); status == stop {
			break
		}
	}
}

func (publisher *Publisher) receive() status {
	select {
	case body := <-publisher.outChannel:
		return publisher.topology(body)
	case <-publisher.signalChannel:
		return stop
	}
}

func (publisher *Publisher) topology(body string) status {
	_, err := publisher.amqpChannel.QueueDeclare("rr", true, false, false,
		false, nil)
	if err != nil {
		return resume
	}
	return publisher.send(body)
}

func (publisher *Publisher) send(body string) status {
	publishing := amqp.Publishing{ContentType: "text/plain", Body: []byte(body)}
	err := publisher.amqpChannel.Publish("", "rr", false, false, publishing)
	if err != nil {
		return stop
	}
	return resume
}
