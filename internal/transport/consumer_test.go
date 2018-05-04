package transport

import (
	"testing"
	"github.com/streadway/amqp"
	"time"
	"taskd/internal/taskd"
)

func TestNewConsumer(t *testing.T) {
	var checkType interface{}
	msg := make(<-chan amqp.Delivery)
	in := make(chan<- taskd.Request)

	c := newConsumer(msg, in)
	t.Run("TestNewConsumer_Return", func(t *testing.T) {
		if c == nil {
			t.Fatal("newConsumer() return nil")
		}
	})
	t.Run("TestNewConsumer_ReturnType", func(t *testing.T) {
		checkType = c
		if _, ok := checkType.(*consumer); !ok {
			t.Fatal("newConsumer() return bad type")
		}
	})
	t.Run("TestNewConsumer_StructInit", func(t *testing.T) {
		if c.in != msg {
			t.Error("c has bad messageChannel")
		}
		if c.outStream != in {
			t.Error("c has bad inChannel")
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

	consumer := consumer{
		in:        messageChannel,
		outStream: inChannel,
	}

	t.Run("TestConsumer_Run_SuccessWay", func(t *testing.T) {
		var request taskd.Request
		go consumer.run()

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
		go consumer.run()

		messageChannel <- delivery
		time.Sleep(time.Second)
		signalChannel <- stopSignal
		time.Sleep(2 * time.Second)
	})

	t.Run("TestConsumer_Run_LockOnReceive", func(t *testing.T) {
		go consumer.run()

		signalChannel <- stopSignal
		time.Sleep(2 * time.Second)
	})
}
