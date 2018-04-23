package rabbitmq

type receiveStatus byte

const (
	stop   receiveStatus = iota
	resume
)
