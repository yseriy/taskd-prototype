package rabbitmq

import (
	"github.com/streadway/amqp"
	"taskd/internal/transport"
)

type receiver struct {
	queue         Queue
	connection    connection
	inputStream   <-chan amqp.Delivery
	toMessage     func(delivery *amqp.Delivery) *transport.Message
	declareParams func(queue *Queue) (string, bool, bool, bool, bool, amqp.Table)
	consumeParams func(queue *Queue) (string, string, bool, bool, bool, bool, amqp.Table)
}

func (r *receiver) Receive() (*transport.Message, error) {
	for {
		delivery, ok := <-r.inputStream
		if !ok {
			if err := r.setupInputStream(); err != nil {
				return nil, err
			}
			continue
		}
		return r.toMessage(&delivery), nil
	}
}

func (r *receiver) setupInputStream() error {
	queue := r.queue
	for {
		channel, err := r.connection.Channel()
		if err != nil {
			return err
		}
		if _, err := channel.QueueDeclare(r.declareParams(&queue)); err != nil {
			if isNotDisconnect(err) {
				return err
			}
			continue
		}
		inputStream, err := channel.Consume(r.consumeParams(&queue))
		if err != nil {
			if isNotDisconnect(err) {
				return err
			}
			continue
		}
		r.inputStream = inputStream
		return nil
	}
}
