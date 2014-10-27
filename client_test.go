package riaken_struct

import (
	"testing"
)

func dial() (*Client, error) {
	marshaller := NewStructMarshal("json", JsonMarshaller, JsonUnmarshaller)
	addrs := []string{"127.0.0.1:10017", "127.0.0.1:10027", "127.0.0.1:10037", "127.0.0.1:1047", "127.0.0.1:1057"}
	client := NewClient(addrs, 3, marshaller)
	err := client.Dial()
	return client, err
}

func TestClient(t *testing.T) {
	client, err := dial()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer client.Close()
}
