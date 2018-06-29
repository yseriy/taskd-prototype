package rabbitmq

type simpleConnection struct {
	connector      func(url string) (amqpConnection, error)
	url            string
	amqpConnection amqpConnection
	amqpChannel    amqpChannel
}

func (c *simpleConnection) Channel() (channel, error) {
	c.closeAll()
	amqpConnection, err := c.connector(c.url)
	if err != nil {
		return nil, err
	}
	amqpChannel, err := amqpConnection.Channel()
	if err != nil {
		return nil, err
	}
	c.amqpConnection, c.amqpChannel = amqpConnection, amqpChannel
	return amqpChannel, nil
}

func (c *simpleConnection) closeAll() {
	if c.amqpChannel != nil {
		c.amqpChannel.Close()
	}
	if c.amqpConnection != nil {
		c.amqpConnection.Close()
	}
	c.amqpChannel, c.amqpConnection = nil, nil
}
