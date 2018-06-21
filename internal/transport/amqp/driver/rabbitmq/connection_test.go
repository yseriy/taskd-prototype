package rabbitmq

import (
	"testing"
	"fmt"
)

func TestConnection_Channel(t *testing.T) {
	//var amqpChannel driver.AMQPChannel
	//var err error

	//t.Run("success way", func(t *testing.T) {
	ec1 := mockRabbitMQChannel{}
	ec2 := mockRabbitMQChannel{}
	if ec1 == ec2 {
		fmt.Println("true")
	}
	fmt.Printf("ec -> %p %p\n", &ec1, &ec2)
	//expectedChannel := channel{rabbitmqChannel: &ec}
	//f1 := func(rabbitmqChannel1 rabbitmqChannel) *channel {
	//	if rabbitmqChannel1 == &ec {
	//		fmt.Printf("%p  %p\n", rabbitmqChannel1, &ec)
	//		return &expectedChannel
	//	}
	//	return nil
	//}
	//mockConnection := mockRabbitMQConnection{
	//	channel: func() (rabbitmqChannel, error) {
	//		m := &mockRabbitMQChannel{}
	//		fmt.Printf("m -> %p\n", m)
	//		return m, nil
	//	},
	//}
	//amqpConnection := connection{rabbitmqConnection: &mockConnection, newChannel: f1}
	//amqpChannel, err = amqpConnection.Channel()
	//if err != nil {
	//	t.Error("amqpConnection.Channel() return error: ", err)
	//}
	//if amqpChannel != &expectedChannel {
	//	t.Error("amqpConnection.Channel() return unexpected value: ", amqpChannel, " ", expectedChannel)
	//}
	//})

	//t.Run("error way", func(t *testing.T) {
	//	expectedError := errors.New("test Channel() error")
	//	mockConnection := mockRabbitMQConnection{
	//		channel: func() (rabbitmqChannel, error) {
	//			return nil, expectedError
	//		},
	//	}
	//	amqpConnection := connection{
	//		rabbitmqConnection: &mockConnection,
	//		mapError: func(err error) error {
	//			return err
	//		},
	//	}
	//	amqpChannel, err = amqpConnection.Channel()
	//	if err != expectedError {
	//		t.Error("amqpConnection.Channel() return unexpected error")
	//	}
	//	if amqpChannel != nil {
	//		t.Error("amqpConnection.Channel() return not nil value")
	//	}
	//})
}

//func TestConnection_Close(t *testing.T) {
//	expectedError := errors.New("test Close() error")
//	cases := []struct {
//		in           rabbitmqConnection
//		expect       error
//		testName     string
//		errorMessage string
//	}{
//		{
//			in:           &mockRabbitMQConnection{},
//			expect:       nil,
//			testName:     "SuccessWay",
//			errorMessage: "amqpConnection.Close() return non nil error",
//		},
//		{
//			in:           &mockRabbitMQConnection{close: func() error { return expectedError },},
//			expect:       expectedError,
//			testName:     "ErrorHappened",
//			errorMessage: "amqpConnection.Close() return unexpected error",
//		},
//	}
//	for _, c := range cases {
//		t.Run(c.testName, func(t *testing.T) {
//			amqpConnection := connection{rabbitmqConnection: c.in}
//			if err := amqpConnection.Close(); err != c.expect {
//				t.Error(c.errorMessage)
//			}
//		})
//	}
//}
