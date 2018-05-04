package connector

import (
	"github.com/streadway/amqp"
	"time"
)

type repeater struct {
	dialer  dialer
	timeout int
}

func newRepeater(dialer dialer, timeout int) dialer {
	return &repeater{dialer: dialer, timeout: timeout}
}

func (repeater *repeater) dial() (*amqp.Connection, error) {
	for {
		connection, err := repeater.dialer.dial()
		if err != nil {
			time.Sleep(time.Duration(repeater.timeout) * time.Second)
			continue
		}
		return connection, nil
	}
}
