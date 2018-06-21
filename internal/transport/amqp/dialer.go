package amqp

import (
	"taskd/internal/transport/amqp/driver"
	"time"
)

type dialer struct {
	connector driver.AMQPConnector
	timeout   time.Duration
}

func (d dialer) Dial() (driver.AMQPConnection, error) {
	for {
		connection, err := d.connector.Dial()
		if err != nil {
			time.Sleep(d.timeout)
			continue
		}
		return connection, nil
	}
}
