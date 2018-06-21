package transport

type (
	Receiver interface {
		Receive() (*Message, error)
	}

	Sender interface {
		Send(message *Message) error
	}

	Message struct {
		Address     string
		ContentType string
		Body        []byte
	}
)
