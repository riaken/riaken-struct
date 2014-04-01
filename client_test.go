package riaken_struct

import (
	"testing"
)

func dial() (*Client, error) {
	marshaller := NewStructMarshal("json", JsonMarshaller, JsonUnmarshaller)
	//addrs := []string{"127.0.0.1:8087"}
	addrs := []string{"127.0.0.1:8083", "127.0.0.1:8084", "127.0.0.1:8085", "127.0.0.1:8086", "127.0.0.1:8087"}
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
