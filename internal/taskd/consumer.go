package taskd

import "taskd/internal/taskd"

type consumer struct {
	queue       amqpQueue
	connector   connector
	converter   converter
	inputStream <-chan taskd.Response
	amqpChannel channel
}

func (c *consumer) Connect(inputStream <-chan taskd.Response) {
	c.inputStream = inputStream
	c.amqpChannel = c.connector.Connect()
	go c.handler()
}

func (c *consumer) handler() {
	for {
		response, ok := <-c.inputStream
		if !ok {
			break
		}
		amqpResponse, err := c.converter.fromTask(response)
		if err != nil {
			continue // drop response
		}
		if err := c.send(amqpResponse); err != nil {
			// drop amqpResponse
		}
	}
}

func (c *consumer) send(response *amqpResponse) error {
	queue := c.queue
	queue.Name = response.RoutingKey

	for {
		if err := c.amqpChannel.QueueDeclare(&queue); err != nil {
			if c.amqpChannel.IsDisconnect(err) {
				c.amqpChannel = c.connector.Connect()
				continue
			} else {
				return err
			}
		}
		if err := c.amqpChannel.Publish(response); err == nil {
			return nil
		}
		c.amqpChannel = c.connector.Connect()
	}
}

func (c *consumer) errorHandle(err error) error {
	if c.amqpChannel.IsDisconnect(err) {
		amqpChannel, err := c.connector.Connect()
		if err != nil {
			panic(err)
		}
		c.amqpChannel = amqpChannel
		return nil
	}
	return err
}
