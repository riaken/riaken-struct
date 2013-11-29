package riaken_struct

import (
	"testing"
)

import (
	"github.com/riaken/riaken-core/rpb"
)

type User struct {
	Name string `json:"name"`
}

func TestObject(t *testing.T) {
	client := dial()
	defer client.Close()
	session, err := client.Session()
	if err != nil {
		t.Error(err.Error())
	}
	defer session.Close()

	user := &User{
		Name: "User",
	}
	bucket := session.GetBucket("b1")
	object := bucket.Object("o1")
	if _, err := object.Store(user); err != nil {
		t.Error(err.Error())
	}

	check := User{}
	if _, err := object.Fetch(&check); err != nil {
		t.Error(err.Error())
	}
	if check.Name != user.Name {
		t.Errorf("got: %s, expected: %s", check.Name, user.Name)
	}

	if _, err := object.Delete(); err != nil {
		t.Error(err.Error())
	}
}

func TestObjectDo(t *testing.T) {
	client := dial()
	defer client.Close()
	session, err := client.Session()
	if err != nil {
		t.Error(err.Error())
	}
	defer session.Close()

	user := &User{
		Name: "User",
	}
	rb := true
	opts := &rpb.RpbPutReq{
		ReturnBody: &rb,
	}
	bucket := session.GetBucket("b1")
	object := bucket.Object("o1")
	if _, err := object.Do(opts).Store(user); err != nil {
		t.Error(err.Error())
	}

	if _, err := object.Delete(); err != nil {
		t.Error(err.Error())
	}
}
