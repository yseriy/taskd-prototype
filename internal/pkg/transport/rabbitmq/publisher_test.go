package rabbitmq

import (
	"testing"
	"github.com/streadway/amqp"
	"taskd/internal/pkg/taskd"
)

func TestNewPublisher(t *testing.T) {
	var checkType interface{}
	in := make(<-chan taskd.Response)
	amqpChannel := &amqp.Channel{}

	p := newPublisher(in, amqpChannel)
	t.Run("TestNewPublisher_Return", func(t *testing.T) {
		if p == nil {
			t.Fatal("newPublisher() return nil")
		}
	})

	t.Run("TestNewPublisher_ReturnType", func(t *testing.T) {
		checkType = p
		if _, ok := checkType.(*publisher); !ok {
			t.Fatal("newPublisher() return bad type")
		}
	})

	t.Run("TestNewPublisher_StructInit", func(t *testing.T) {
		if p.outAMQPStream != amqpChannel {
			t.Error("Publishet has bad outAMQPStream")
		}
		if p.inStream != in {
			t.Error("Publishet has bad inStream")
		}
	})
}
