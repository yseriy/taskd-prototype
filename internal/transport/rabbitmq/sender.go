package rabbitmq

import (
	"github.com/streadway/amqp"
	"taskd/internal/transport"
)

type sender struct {
	queue         Queue
	connection    connection
	channel       channel
	declareParams func(queue *Queue) (string, bool, bool, bool, bool, amqp.Table)
	publishParams func(message *transport.Message) (string, string, bool, bool, amqp.Publishing)
}

func (s *sender) Send(message *transport.Message) error {
	queue := s.initQueue(message)
	for {
		if _, err := s.channel.QueueDeclare(s.declareParams(queue)); err != nil {
			if err := s.reconnectIfDisconnect(err); err != nil {
				return err
			}
			continue
		}
		if err := s.channel.Publish(s.publishParams(message)); err != nil {
			if err := s.reconnectIfDisconnect(err); err != nil {
				return err
			}
			continue
		}
		return nil
	}
}

func (s *sender) initQueue(message *transport.Message) *Queue {
	queue := s.queue
	queue.Name = message.Address
	return &queue
}

func (s *sender) reconnectIfDisconnect(err error) error {
	if isDisconnect(err) {
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
