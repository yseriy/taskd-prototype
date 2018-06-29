package rabbitmq

import (
	"testing"
	"github.com/golang/mock/gomock"
	"taskd/internal/transport"
	"github.com/streadway/amqp"
)

func TestSender_Send(t *testing.T) {
	testQueue := Queue{}
	testMessage := transport.Message{}
	testPublishing := amqp.Publishing{
		ContentType: testMessage.ContentType,
		Body:        testMessage.Body,
	}
	declareParams := func(queue *Queue) (string, bool, bool, bool, bool, amqp.Table) {
		return queue.Name, queue.Durable, queue.AutoDelete, queue.Exclusive, false, amqp.Table(queue.Args)
	}
	publishParams := func(message *transport.Message) (string, string, bool, bool, amqp.Publishing) {
		return "", message.Address, false, true, amqp.Publishing{
			ContentType: message.ContentType,
			Body:        message.Body,
		}
	}

	setupSuccessWay := func(controller *gomock.Controller) (channel, connection) {
		mockChannel := NewMockchannel(controller)
		mockConnection := NewMockconnection(controller)
		mockConnection.EXPECT().Channel().Return(mockChannel, nil).Times(0)
		gomock.InOrder(
			mockChannel.EXPECT().QueueDeclare(testMessage.Address, testQueue.Durable, testQueue.AutoDelete,
				testQueue.Exclusive, false, amqp.Table(testQueue.Args)).Return(nil, nil).Times(1),
			mockChannel.EXPECT().Publish("", testMessage.Address, true, false,
				testPublishing).Return(nil).Times(1),
		)
		return mockChannel, mockConnection
	}

	cases := []struct {
		setupMock     func(controller *gomock.Controller) (channel, connection)
		expectedError error
		testName      string
	}{
		{setupMock: setupSuccessWay, expectedError: nil, testName: "SuccessWay"},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockChannel, mockConnection := c.setupMock(mockCtrl)
			sender := &sender{
				queue:         testQueue,
				connection:    mockConnection,
				channel:       mockChannel,
				declareParams: declareParams,
				publishParams: publishParams,
			}

			err := sender.Send(&testMessage)
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
