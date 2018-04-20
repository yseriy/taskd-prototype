package taskd

import (
	"github.com/streadway/amqp"
	"fmt"
)

type status byte

const (
	stop   = 0
	resume = 1
)

type Consumer struct {
	messageChannel <-chan amqp.Delivery
	inChannel      chan<- string
	signalChannel  <-chan struct{}
}

func NewConsumer(message <-chan amqp.Delivery, in chan<- string, signal <-chan struct{}) *Consumer {
	return &Consumer{messageChannel: message, inChannel: in, signalChannel: signal}
}

func (consumer *Consumer) Run() {
	go consumer.handle()
}

func (consumer *Consumer) handle() {
	fmt.Println("start receiving")
	for {
		if status := consumer.receive(); status == stop {
			fmt.Println("stop receiving")
			break
		}
		fmt.Println("processing complete")
	}
}

func (consumer *Consumer) receive() status {
	select {
	case delivery := <-consumer.messageChannel:
		fmt.Println("start processing new delivery")
		return consumer.convert(&delivery)
	case <-consumer.signalChannel:
		fmt.Println("interupt in receive")
		return stop
	}
}

func (consumer *Consumer) convert(delivery *amqp.Delivery) status {
	body := string(delivery.Body)
	return consumer.send(body)
}

func (consumer *Consumer) send(body string) status {
	fmt.Println("send body")
	select {
	case consumer.inChannel <- body:
		return resume
	case <-consumer.signalChannel:
		fmt.Println("interupt in send")
		return stop
	}
}
