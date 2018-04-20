package taskd

import (
	"github.com/streadway/amqp"
	"fmt"
	"time"
)

func worker(in chan amqp.Delivery, out chan string) {
	for {
		fmt.Println("start")
		delivery := <-in
		time.Sleep(5000 * time.Millisecond)
		out <- string(delivery.Body)
		fmt.Println("stop")
	}
}
