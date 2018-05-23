package rabbitmq

import (
	"github.com/streadway/amqp"
	"testing"
	"taskd/internal/taskd"
	"fmt"
)

type mockConnector struct {
	connectCount int
	input        chan amqp.Delivery
}

func NewMockConnector(input chan amqp.Delivery) *mockConnector {
	return &mockConnector{connectCount: 0, input: input}
}

func (connector *mockConnector) Connect() channel {
	connector.connectCount++
	if connector.connectCount > 3 {
		panic("stop test")
	}
	return &mockChannel{input: connector.input}
}

type mockChannel struct {
	input chan amqp.Delivery
}

func (channel *mockChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

func (channel *mockChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return channel.input, nil
}

func (channel *mockChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return nil
}

func TestConsumer(t *testing.T) {
	input := make(chan amqp.Delivery)
	output := make(chan taskd.Request)
	connector := NewMockConnector(input)

	go func() {
		defer fmt.Println("P1 stop")
		for i := 0; i < 3; i++ {
			input <- amqp.Delivery{}
			fmt.Println("P1")
		}
		close(input)
	}()

	go func() {
		defer fmt.Println("P2 stop")
		for {
			o := <-output
			fmt.Println("P2 ", o)
		}
	}()

	consumer := newConsumer(output, connector, newConverter(), MainQueue{}, ConsumeParams{})
	consumer.run()
}
