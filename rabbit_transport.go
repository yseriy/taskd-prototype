package taskd

import (
	"github.com/streadway/amqp"
	"fmt"
	"time"
)

type RabbitTransport struct {
	rabbitURI string
}

func NewRabbitTransport(rabbitURI string) *RabbitTransport {
	return &RabbitTransport{rabbitURI: rabbitURI}
}

func (transport *RabbitTransport) Start(in chan<- string, out <-chan string) error {
	timeout := 5
	stop := make(chan struct{})
	for {
		connection, err := amqp.Dial(transport.rabbitURI)
		if err != nil {
			fmt.Println("error: ", err, ". try restart connection after 5 sec.")
			time.Sleep(time.Duration(timeout) * time.Second)
			continue
		}
		if err := start(connection, in, out, stop); err != nil {
			connection.Close()
			return err
		}
		fmt.Println("connected")
		notify := connection.NotifyClose(make(chan *amqp.Error))
		err = <-notify
		fmt.Println("error: ", err, ". try restart connection")
	}
	return nil
}

func start(connection *amqp.Connection, in chan<- string, out <-chan string, stop <-chan struct{}) error {
	//var err error
	//errorFlag := false
	//
	//if err = startConsumer(connection); err != nil {
	//	errorFlag = true
	//}
	//if err = startPublisher(connection); err != nil {
	//	errorFlag = true
	//}
	//if errorFlag {
	//	fmt.Println("error: ", err, ". try restart connection after 5 sec.")
	//	return err
	//}
	return nil
}

func startConsumer(connection *amqp.Connection, in chan<- string, stop <-chan struct{}) error {
	consumerChannel, err := connection.Channel()
	if err != nil {
		return err
	}
	queue, err := consumerChannel.QueueDeclare("test_queue_122", true, false,
		false, false, nil)
	if err != nil {
		return err
	}
	msg, err := consumerChannel.Consume(queue.Name, "", true, false,
		false, false, nil)
	if err != nil {
		return err
	}
	go consumer(msg, in, stop)
	return nil
}

func consumer(msg <-chan amqp.Delivery, in chan<- string, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			break
		case d := <-msg:
			if sendTo(in, stop, d.Body) {
				break
			}
		}
	}
}

func sendTo(in chan<- string, stop <-chan struct{}, body []byte) bool {
	select {
	case in <- string(body):
	case <-stop:
		return false
	}
	return true
}

func startPublisher(connection *amqp.Connection) error {
	publisherChannel, err := connection.Channel()
	if err != nil {
		return err
	}
	publisherChannel.Close()
	return nil
}

func publisher(channel *amqp.Channel, in chan string, notify chan struct{}) {
	for {
		select {
		case <-notify:
			break
		case body := <-in:
			publishing := amqp.Publishing{ContentType: "text/plain", Body: []byte(body)}
			err := channel.Publish("", "rr", false, false, publishing)
			if err != nil {
				<-notify
				break
			}
		}
	}
}

//func waitCloseNotify(notify <-chan *amqp.Error) {
//	fmt.Println("connected")
//	//if err, ok := <-notify; !ok {
//	err := <-notify
//	fmt.Println("error: ", err, ". try restart connection")
//	//} else {
//	//	fmt.Println("disconnect. try restart connection")
//	//}
//}
