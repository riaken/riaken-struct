package riaken_struct

import (
	"testing"
)

func TestSession(t *testing.T) {
	client := dial()
	defer client.Close()
	session := client.Session()
	defer session.Release()
}
