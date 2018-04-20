package taskd

import (
	"testing"
	"github.com/streadway/amqp"
	"time"
)

func TestNewConsumer(t *testing.T) {
	var checkType interface{}
	msg := make(<-chan amqp.Delivery)
	in := make(chan<- string)
	stop := make(<-chan struct{})

	consumer := NewConsumer(msg, in, stop)
	t.Run("TestNewConsumer_Return", func(t *testing.T) {
		if consumer == nil {
			t.Fatal("NewConsumer() return nil")
		}
	})
	t.Run("TestNewConsumer_ReturnType", func(t *testing.T) {
		checkType = consumer
		if _, ok := checkType.(*Consumer); !ok {
			t.Fatal("NewConsumer() return bad type")
		}
	})
	t.Run("TestNewConsumer_StructInit", func(t *testing.T) {
		if consumer.messageChannel != msg {
			t.Error("Consumer has bad messageChannel")
		}
		if consumer.inChannel != in {
			t.Error("Consumer has bad inChannel")
		}
		if consumer.signalChannel != stop {
			t.Error("Consumer has bad signalChannel")
		}
	})
}

func TestConsumer_Run(t *testing.T) {
	messageChannel := make(chan amqp.Delivery)
	inChannel := make(chan string)
	signalChannel := make(chan struct{})

	testBody := "test_body"
	stopSignal := struct{}{}
	delivery := amqp.Delivery{Body: []byte(testBody)}

	consumer := Consumer{
		messageChannel: messageChannel,
		inChannel:      inChannel,
		signalChannel:  signalChannel,
	}

	t.Run("TestConsumer_Run_SuccessWay", func(t *testing.T) {
		var body string
		consumer.Run()

		for i := 0; i < 5; i++ {
			messageChannel <- delivery
			body = <-inChannel
		}
		signalChannel <- stopSignal

		if body != testBody {
			t.Error("For message body expected '", testBody, "got '", body, "'")
		}
	})

	t.Run("TestConsumer_Run_LockOnSend", func(t *testing.T) {
		consumer.Run()

		messageChannel <- delivery
		time.Sleep(time.Second)
		signalChannel <- stopSignal
		time.Sleep(2 * time.Second)
	})

	t.Run("TestConsumer_Run_LockOnReceive", func(t *testing.T) {
		consumer.Run()

		signalChannel <- stopSignal
		time.Sleep(2 * time.Second)
	})
}
