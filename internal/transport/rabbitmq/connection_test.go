package rabbitmq

import (
	"testing"
	"github.com/golang/mock/gomock"
	"errors"
)

func TestSimpleConnection_Channel(t *testing.T) {
	setupSuccessWayConnect := func(controller *gomock.Controller) (*simpleConnection, amqpChannel,
		amqpConnection, error) {
		mockUrl := "test url 1096100260"
		mockChannel := NewMockamqpChannel(controller)
		mockConnection := NewMockamqpConnection(controller)
		mockConnector := func(url string) (amqpConnection, error) {
			if url == mockUrl {
				return mockConnection, nil
			}
			return nil, errors.New("unexpected url")
		}
		connection := &simpleConnection{
			connector: mockConnector,
			url:       mockUrl,
		}
		mockConnection.EXPECT().Channel().Return(mockChannel, nil).Times(1)
		return connection, mockChannel, mockConnection, nil
	}

	setupSuccessWayReconnect := func(controller *gomock.Controller) (*simpleConnection, amqpChannel,
		amqpConnection, error) {
		mockOldChannel := NewMockamqpChannel(controller)
		mockOldConnection := NewMockamqpConnection(controller)
		mockUrl := "test url 1036105260"
		mockChannel := NewMockamqpChannel(controller)
		mockConnection := NewMockamqpConnection(controller)
		mockConnector := func(url string) (amqpConnection, error) {
			if url == mockUrl {
				return mockConnection, nil
			}
			return nil, errors.New("unexpected url")
		}
		connection := &simpleConnection{
			connector:      mockConnector,
			url:            mockUrl,
			amqpConnection: mockOldConnection,
			amqpChannel:    mockOldChannel,
		}
		gomock.InOrder(
			mockOldChannel.EXPECT().Close().Return(nil).Times(1),
			mockOldConnection.EXPECT().Close().Return(nil).Times(1),
			mockConnection.EXPECT().Channel().Return(mockChannel, nil).Times(1),
		)
		return connection, mockChannel, mockConnection, nil
	}

	setupConnectorErrorConnect := func(controller *gomock.Controller) (*simpleConnection, amqpChannel,
		amqpConnection, error) {
		expectedUrl := "test url 5017194055"
		expectedError := errors.New("test connector error 5017194055")
		mockConnector := func(url string) (amqpConnection, error) {
			if url == expectedUrl {
				return nil, expectedError
			}
			return nil, errors.New("unexpected url")
		}
		connection := &simpleConnection{
			connector: mockConnector,
			url:       expectedUrl,
		}
		return connection, nil, nil, expectedError
	}

	setupConnectorErrorReconnect := func(controller *gomock.Controller) (*simpleConnection, amqpChannel,
		amqpConnection, error) {
		mockOldChannel := NewMockamqpChannel(controller)
		mockOldConnection := NewMockamqpConnection(controller)
		expectedUrl := "test url 3388713263"
		expectedError := errors.New("test connector error 3388713263")
		mockConnector := func(url string) (amqpConnection, error) {
			if url == expectedUrl {
				return nil, expectedError
			}
			return nil, errors.New("unexpected url")
		}
		connection := &simpleConnection{
			connector:      mockConnector,
			url:            expectedUrl,
			amqpConnection: mockOldConnection,
			amqpChannel:    mockOldChannel,
		}
		gomock.InOrder(
			mockOldChannel.EXPECT().Close().Return(nil).Times(1),
			mockOldConnection.EXPECT().Close().Return(nil).Times(1),
		)
		return connection, nil, nil, expectedError
	}

	setupConnectionErrorConnect := func(controller *gomock.Controller) (*simpleConnection, amqpChannel,
		amqpConnection, error) {
		expectedUrl := "test url 2139529643"
		expectedError := errors.New("test connection error 2139529643")
		mockConnection := NewMockamqpConnection(controller)
		mockConnector := func(url string) (amqpConnection, error) {
			if url == expectedUrl {
				return mockConnection, nil
			}
			return nil, errors.New("unexpected url")
		}
		connection := &simpleConnection{
			connector: mockConnector,
			url:       expectedUrl,
		}
		mockConnection.EXPECT().Channel().Return(nil, expectedError).Times(1)
		return connection, nil, nil, expectedError
	}

	setupConnectionErrorReconnect := func(controller *gomock.Controller) (*simpleConnection, amqpChannel,
		amqpConnection, error) {
		mockOldChannel := NewMockamqpChannel(controller)
		mockOldConnection := NewMockamqpConnection(controller)
		expectedUrl := "test url 9239150110"
		expectedError := errors.New("test connection error 9239150110")
		mockConnection := NewMockamqpConnection(controller)
		mockConnector := func(url string) (amqpConnection, error) {
			if url == expectedUrl {
				return mockConnection, nil
			}
			return nil, errors.New("unexpected url")
		}
		connection := &simpleConnection{
			connector:      mockConnector,
			url:            expectedUrl,
			amqpConnection: mockOldConnection,
			amqpChannel:    mockOldChannel,
		}
		gomock.InOrder(
			mockOldChannel.EXPECT().Close().Return(nil).Times(1),
			mockOldConnection.EXPECT().Close().Return(nil).Times(1),
			mockConnection.EXPECT().Channel().Return(nil, expectedError).Times(1),
		)
		return connection, nil, nil, expectedError
	}

	cases := []struct {
		setupFunc func(controller *gomock.Controller) (*simpleConnection, amqpChannel, amqpConnection, error)
		testName  string
	}{
		{setupFunc: setupSuccessWayConnect, testName: "SuccessWayConnect"},
		{setupFunc: setupSuccessWayReconnect, testName: "SuccessWayReconnect"},
		{setupFunc: setupConnectorErrorConnect, testName: "ConnectorErrorConnect"},
		{setupFunc: setupConnectorErrorReconnect, testName: "ConnectorErrorReconnect"},
		{setupFunc: setupConnectionErrorConnect, testName: "ConnectionErrorConnect"},
		{setupFunc: setupConnectionErrorReconnect, testName: "ConnectionErrorReconnect"},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			connection, mockChannel, mockConnection, expectedError := c.setupFunc(mockCtrl)

			testChannel, err := connection.Channel()
			if err != expectedError {
				t.Error("connection.Channel() return non nil error: ", err)
			}
			if testChannel != mockChannel {
				t.Error("connection.Channel() return non expected value")
			}
			if connection.amqpChannel != mockChannel {
				t.Error("connection.amqpChannel has non expected value")
			}
			if connection.amqpConnection != mockConnection {
				t.Error("connection.amqpConnection has non expected value")
			}
		})
	}
}
