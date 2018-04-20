package taskd

import "testing"

func TestRabbitTransport_NewRabbitTransport(t *testing.T) {
	expectedURI := "test_uri"
	var checkType interface{}

	transport := NewRabbitTransport(expectedURI)
	if transport == nil {
		t.Fatal("NewRabbitTransport() return nil")
	}

	checkType = transport
	if _, ok := checkType.(*RabbitTransport); !ok {
		t.Fatal("Object type is not '*RabbitTranport'")
	}

	if expectedURI != transport.rabbitURI {
		t.Error("Cannot set rabbit uri")
	}
}

func TestRabbitTransport_Start(t *testing.T) {
	in, out := make(chan string), make(chan string)
	transport := RabbitTransport{rabbitURI: "amqp://devel:devel@ds:5672/"}
	if err := transport.Start(in, out); err != nil {
		t.Error("transport.Start() return '", err, "'")
	}
}

func TestRabbitTransport_Restart(t *testing.T) {

}
