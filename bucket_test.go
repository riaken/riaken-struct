package riaken_struct

import (
	"testing"
)

func TestBucket(t *testing.T) {
	client, err := dial()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer client.Close()
	session := client.Session()
	defer session.Release()

	session.GetBucket("b1")
}
