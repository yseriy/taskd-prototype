package transport

type connector interface {
	Connect() (connection, error)
}

type connection interface {
	Channel() (channel, error)
	Close() error
}

type channel interface {
	QueueDeclare(queue amqpQueue) error
	Publish(response amqpResponse) error
}

type amqpQueue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       table
}

type table map[string]interface{}

type amqpResponse struct {
	Exchange   string
	RoutingKey string
	Body       []byte
	Mandatory  bool
	Immediate  bool
}
