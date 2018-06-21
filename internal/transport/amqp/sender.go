package amqp

import (
	"taskd/internal/transport"
	"taskd/internal/transport/amqp/driver"
	"time"
)

type sender struct {
	queue      driver.AMQPQueue
	connection driver.AMQPConnection
	channel    driver.AMQPChannel
}

func NewSender(url string, timeout time.Duration, queue driver.AMQPQueue) (*sender, error) {
	connection, err := newManager(url, timeout)
	if err != nil {
		return nil, err
	}
	sender := &sender{queue: queue, connection: connection}
	if err := sender.setupChannel(); err != nil {
		return nil, err
	}
	return sender, nil
}

func (s *sender) Send(message *transport.Message) error {
	queue := s.initQueue(message)
	for {
		if err := s.channel.QueueDeclare(queue); err != nil {
			if err := s.reconnectIfDisconnect(err); err != nil {
				return err
			} else {
				continue
			}
		}
		if err := s.channel.Publish(message); err != nil {
			if err := s.reconnectIfDisconnect(err); err != nil {
				return err
			} else {
				continue
			}
		}
		return nil
	}
}

func (s *sender) initQueue(message *transport.Message) *driver.AMQPQueue {
	queue := s.queue
	queue.Name = message.Address
	return &queue
}

func (s *sender) reconnectIfDisconnect(err error) error {
	if driver.IsDisconnect(err) {
		return s.setupChannel()
	}
	return err
}

func (s *sender) setupChannel() error {
	channel, err := s.connection.Channel()
	if err != nil {
		return err
	}
	s.channel = channel
	return nil
}
