package transport

import "time"

type dialer interface {
	Dial() connection
}

type simpleDialer struct {
	timeout   time.Duration
	connector connector
}

func (d simpleDialer) Dial() connection {
	for {
		connection, err := d.connector.Connect()
		if err != nil {
			time.Sleep(d.timeout)
			continue
		}
		return connection
	}
}
