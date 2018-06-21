package amqp

import (
	"taskd/internal/transport"
	"taskd/internal/transport/amqp/driver"
	"time"
)

type receiver struct {
	queue       driver.AMQPQueue
	connection  driver.AMQPConnection
	inputStream <-chan transport.Message
}

func NewReceiver(url string, timeout time.Duration, queue driver.AMQPQueue) (*receiver, error) {
	connection, err := newManager(url, timeout)
	if err != nil {
		return nil, err
	}
	receiver := &receiver{queue: queue, connection: connection}
	if err := receiver.setupInputStream(); err != nil {
		return nil, err
	}
	return receiver, nil
}

func (r *receiver) Receive() (*transport.Message, error) {
	for {
		amqpRequest, ok := <-r.inputStream
		if !ok {
			if err := r.setupInputStream(); err != nil {
				return nil, err
			} else {
				continue
			}
		}
		return &amqpRequest, nil
	}
}

func (r *receiver) setupInputStream() error {
	queue := r.queue
	for {
		channel, err := r.connection.Channel()
		if err != nil {
			return err
		}
		if err := channel.QueueDeclare(&queue); err != nil {
			if driver.IsDisconnect(err) {
				continue
			} else {
				return err
			}
		}
		inputStream, err := channel.Consume(&queue)
		if err != nil {
			if driver.IsDisconnect(err) {
				continue
			} else {
				return err
			}
		}
		r.inputStream = inputStream
		return nil
	}
}
