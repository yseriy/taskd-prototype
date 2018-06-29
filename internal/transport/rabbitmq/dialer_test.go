package rabbitmq

import (
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestDialer_Channel(t *testing.T) {
	setupSuccessWay := func(controller *gomock.Controller) (*dialer, channel, time.Duration, error) {
		timeout := time.Millisecond
		mockChannel := NewMockchannel(controller)
		mockConnection := NewMockconnection(controller)
		dialer := &dialer{connection: mockConnection, timeout: timeout}
		mockConnection.EXPECT().Channel().Return(mockChannel, nil).Times(1)
		return dialer, mockChannel, timeout, nil
	}

	setupConnectionError3Times := func(controller *gomock.Controller) (*dialer, channel, time.Duration, error) {
		timeout := time.Millisecond
		mockChannel := NewMockchannel(controller)
		mockConnection := NewMockconnection(controller)
		dialer := &dialer{connection: mockConnection, timeout: timeout}
		gomock.InOrder(
			mockConnection.EXPECT().Channel().Return(nil, errors.New("test connection error 1")).Times(1),
			mockConnection.EXPECT().Channel().Return(nil, errors.New("test connection error 2")).Times(1),
			mockConnection.EXPECT().Channel().Return(nil, errors.New("test connection error 3")).Times(1),
			mockConnection.EXPECT().Channel().Return(mockChannel, nil).Times(1),
		)
		return dialer, mockChannel, timeout, nil
	}

	cases := []struct {
		setupFunc func(controller *gomock.Controller) (*dialer, channel, time.Duration, error)
		testName  string
	}{
		{setupFunc: setupSuccessWay, testName: "SuccessWay"},
		{setupFunc: setupConnectionError3Times, testName: "ConnectionError3Times"},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			dialer, mockChannel, expectedTimeout, expectedError := c.setupFunc(mockCtrl)

			testChannel, err := dialer.Channel()
			if err != expectedError {
				t.Error("dialer.Channel() return non nil error: ", err)
			}
			if testChannel != mockChannel {
				t.Error("dialer.Channel() return non expected value")
			}
			if dialer.timeout != expectedTimeout {
				t.Error("dialer.timeout has unexpected value")
			}
		})
	}
}
