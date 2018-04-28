package main

import (
	"github.com/streadway/amqp"
	"log"
	"fmt"
)

func main() {
	conn, err := amqp.Dial("amqp://devel:devel@ds:5672/")
	if err != nil {
		log.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	q, err := ch.QueueDeclare("test_queue_log", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := ch.QueueBind(q.Name, "#", "amq.rabbitmq.trace", false, nil); err != nil {
		log.Fatal(err)
	}
	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	for m := range msg {
		fmt.Println(m.Headers, string(m.Body))
	}
	conn.Close()
}
