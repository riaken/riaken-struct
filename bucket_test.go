package riaken_struct

import (
	"testing"
)

func TestBucket(t *testing.T) {
	client := dial()
	defer client.Close()
	session, err := client.Session()
	if err != nil {
		t.Error(err.Error())
	}
	defer session.Close()

	session.GetBucket("b1")
}
