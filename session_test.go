package riaken_struct

import (
	"testing"
)

func TestSession(t *testing.T) {
	client, err := dial()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer client.Close()
	session := client.Session()
	defer session.Release()
}
