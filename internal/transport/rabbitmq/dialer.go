package rabbitmq

import "time"

type dialer struct {
	connection connection
	timeout    time.Duration
}

func (d dialer) Channel() (channel, error) {
	for {
		channel, err := d.connection.Channel()
		if err != nil {
			time.Sleep(d.timeout)
			continue
		}
		return channel, nil
	}
}
