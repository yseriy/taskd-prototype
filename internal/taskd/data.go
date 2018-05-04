package taskd

type Request struct {
	Id      string
	TaskId  string
	Account string
	Address string
	Body    map[string]string
}

type Response struct {
	Id      string
	TaskId  string
	Account string
	Address string
	Body    map[string]string
}
