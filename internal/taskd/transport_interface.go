package taskd

type Transport interface {
	Connect(chan<- Request, <-chan Response)
}
