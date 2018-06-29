package rabbitmq

import (
	"testing"
	"github.com/golang/mock/gomock"
	"taskd/internal/transport"
	"github.com/streadway/amqp"
)

func TestSender_Send(t *testing.T) {
	setupSuccessWay := func(controller *gomock.Controller, queue Queue,
		message *transport.Message) (channel, connection) {
		mockChannel := NewMockchannel(controller)
		mockConnection := NewMockconnection(controller)
		publishing := amqp.Publishing{
			ContentType: message.ContentType,
			Body:        message.Body,
		}

		mockConnection.EXPECT().Channel().Return(mockChannel, nil).Times(1)
		gomock.InOrder(
			mockChannel.EXPECT().QueueDeclare(message.Address, queue.Durable, queue.AutoDelete, queue.Exclusive,
				false, amqp.Table(queue.Args)).Return(nil, nil).Times(1),
			mockChannel.EXPECT().Publish("", message.Address, true, false,
				publishing).Return(nil).Times(1),
		)

		return mockChannel, mockConnection
	}

	cases := []struct {
		setupMock     func(controller *gomock.Controller, queue Queue) (channel, connection)
		inputMessage  *transport.Message
		expectedError error
		testName      string
	}{
		{setupMock: setupSuccessWay, testName: "SuccessWay"},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockChannel, mockConnection := c.setupMock(mockCtrl, )
			sender := &sender{}

			err := sender.Send(c.inputMessage)
			if err != c.expectedError {
				t.Error("sender.Send() return unexpected error")
			}
			if sender.channel != mockChannel {
				t.Error("sender.channel has unexpected value")
			}
			if sender.connection != mockConnection {
				t.Error("sender.connection has unexpected value")
			}
		})
	}
}
