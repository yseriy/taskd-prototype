package main

import (
	"github.com/streadway/amqp"
	"log"
	"time"
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

	//notify := ch.NotifyClose(make(chan *amqp.Error))
	//go func() {
	//	for {
	//		err, ok := <-notify
	//		if !ok {
	//			break
	//		}
	//		log.Print("p3 ", err)
	//	}
	//}()

	confirm := ch.NotifyPublish(make(chan amqp.Confirmation))
	go func() {
		for {
			err, ok := <-confirm
			if !ok {
				break
			}
			log.Print("p7 ", err.DeliveryTag, "---------------->", err.Ack)
		}
	}()

	if err := ch.Confirm(false); err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	//log.Print("start sleep")
	//time.Sleep(30 * time.Second)
	//log.Print("stop sleep")

	//_, err = ch.QueueDeclare("test_queue", true, false, false, false, nil)
	//if err != nil {
	//	log.Print("p1 ", err)
	//}

	//time.Sleep(2 * time.Second)

	//ch, err = conn.Channel()
	//if err != nil {
	//	log.Fatal("p6 ", err)
	//}

	//msg, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	//if err != nil {
	//	log.Fatal("p2 ", err)
	//}
	//
	//<-msg

	for i := 0; i < 9; i++ {
		err = ch.Publish("", "testr", true, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("test_msg"),
		})
		if err != nil {
			log.Fatal("p2 ", err)
		}
	}

	time.Sleep(2 * time.Second)
	ch.Close()
	conn.Close()
	time.Sleep(2 * time.Second)
}
