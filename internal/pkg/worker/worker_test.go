package worker

import (
	"testing"
	"github.com/streadway/amqp"
	"time"
)

func TestWorker(t *testing.T) {
	s := "test_string"
	delivery := amqp.Delivery{Body: []byte(s)}

	in := make(chan amqp.Delivery)
	out := make(chan string)

	a := 3

	for i := 0; i < a; i++ {
		go worker(in, out)
	}

	for i := 0; i < a; i++ {
		in <- delivery
	}

	for i := 0; i < a; i++ {
		_ = <-out
	}
	//for s1 := range out {
	//	fmt.Println(s1)
	//}
	time.Sleep(120000 * time.Millisecond)
	//if s != s1 {
	//	t.Error("Expected string :'", s, "' got strung: '", s1)
	//}
}
