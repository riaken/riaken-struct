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

	b := session.GetBucket("b1")
	b.Object("o1")
	b.Counter("c1")
	b.Crdt("dt1")
}
