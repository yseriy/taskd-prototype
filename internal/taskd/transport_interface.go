package taskd

type transport interface {
	Connect(chan<- Request, <-chan Response)
}
