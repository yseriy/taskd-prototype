package rabbitmq

import (
	"testing"
	"github.com/streadway/amqp"
	"time"
	"taskd"
)

func TestNewConsumer(t *testing.T) {
	var checkType interface{}
	msg := make(<-chan amqp.Delivery)
	in := make(chan<- taskd.Request)
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
		if consumer.in != msg {
			t.Error("Consumer has bad messageChannel")
		}
		if consumer.out != in {
			t.Error("Consumer has bad inChannel")
		}
		if consumer.signal != stop {
			t.Error("Consumer has bad signal")
		}
	})
}

func TestConsumer_Run(t *testing.T) {
	messageChannel := make(chan amqp.Delivery)
	inChannel := make(chan taskd.Request)
	signalChannel := make(chan struct{})

	testBody := "test_body"
	stopSignal := struct{}{}
	delivery := amqp.Delivery{Body: []byte(testBody)}

	consumer := Consumer{
		in:     messageChannel,
		out:    inChannel,
		signal: signalChannel,
	}

	t.Run("TestConsumer_Run_SuccessWay", func(t *testing.T) {
		var request taskd.Request
		go consumer.Run()

		for i := 0; i < 5; i++ {
			messageChannel <- delivery
			request = <-inChannel
		}
		signalChannel <- stopSignal

		if request.Body != testBody {
			t.Error("For message request expected '", testBody, "got '", request, "'")
		}
	})

	t.Run("TestConsumer_Run_LockOnSend", func(t *testing.T) {
		go consumer.Run()

		messageChannel <- delivery
		time.Sleep(time.Second)
		signalChannel <- stopSignal
		time.Sleep(2 * time.Second)
	})

	t.Run("TestConsumer_Run_LockOnReceive", func(t *testing.T) {
		go consumer.Run()

		signalChannel <- stopSignal
		time.Sleep(2 * time.Second)
	})
}
