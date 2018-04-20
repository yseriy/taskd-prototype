package taskd

import (
	"testing"
	"github.com/streadway/amqp"
)

func TestNewPublisher(t *testing.T) {
	var checkType interface{}
	amqpChannel := &amqp.Channel{}
	out := make(<-chan string)
	stop := make(<-chan struct{})

	publisher := NewPublisher(amqpChannel, out, stop)
	t.Run("TestNewPublisher_Return", func(t *testing.T) {
		if publisher == nil {
			t.Fatal("NewPublisher() return nil")
		}
	})

	t.Run("TestNewPublisher_ReturnType", func(t *testing.T) {
		checkType = publisher
		if _, ok := checkType.(*Publisher); !ok {
			t.Fatal("NewPublisher() return bad type")
		}
	})

	t.Run("TestNewPublisher_StructInit", func(t *testing.T) {
		if publisher.amqpChannel != amqpChannel {
			t.Error("Publishet has bad amqpChannel")
		}
		if publisher.outChannel != out {
			t.Error("Publishet has bad outChannel")
		}
		if publisher.signalChannel != stop {
			t.Error("Publishet has bad signalChannel")
		}
	})
}
