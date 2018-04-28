package taskd

type Request struct {
	Body           string
	TransportStash interface{}
}

type Response struct {
	Body           string
	TransportStash interface{}
}
