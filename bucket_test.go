package riaken_struct

import (
	"testing"
)

func TestBucket(t *testing.T) {
	client := dial()
	defer client.Close()
	session := client.Session()
	defer session.Release()

	session.GetBucket("b1")
}
