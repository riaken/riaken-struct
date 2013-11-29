package riaken_struct

import (
	"testing"
	"time"
)

func dial() *Client {
	marshaller := NewStructMarshal("json", JsonMarshaller, JsonUnmarshaller)
	addrs := []string{"127.0.0.1:8087"}
	//addrs := []string{"127.0.0.1:8083", "127.0.0.1:8084", "127.0.0.1:8085", "127.0.0.1:8086", "127.0.0.1:8087"}
	client := NewClient(addrs, 3, time.Second*2, marshaller)
	//client.debug = true
	client.Dial()
	return client
}

func TestClient(t *testing.T) {
	client := dial()
	defer client.Close()
}
