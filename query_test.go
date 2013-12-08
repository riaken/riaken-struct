package riaken_struct

import (
	"os/exec"
	"testing"
)

type Animal struct {
	Id   string
	Name string `json:"name" riak:"index"`
}

func TestQuery(t *testing.T) {
	client := dial()
	defer client.Close()
	session := client.Session()
	defer session.Release()

	session.Query()
}

func TestSecondaryIndexes(t *testing.T) {
	client := dial()
	defer client.Close()
	session := client.Session()
	defer session.Release()

	// Setup
	animal := Animal{
		Name: "chicken",
	}
	bucket := session.GetBucket("animals")
	object := bucket.Object("2i-a1")
	if _, err := object.Store(&animal); err != nil {
		t.Error(err.Error())
	}

	// Query
	var check []Animal
	query := session.Query()
	if _, err := query.Key("Id").Out(&check).SecondaryIndexes([]byte("animals"), []byte("name_bin"), []byte("chicken"), nil, nil, 0, nil); err != nil {
		t.Error(err.Error())
	}

	if len(check) == 0 {
		t.Error("expected results")
	} else {
		if check[0].Name != "chicken" {
			t.Errorf("expected: chicken, got: %s", check[0].Name)
		}
		if check[0].Id != "2i-a1" {
			t.Errorf("expected: 2i-a1, got: %s", check[0].Id)
		}
	}

	// Cleanup
	if _, err := object.Delete(); err != nil {
		t.Error(err.Error())
	}
}

func TestSearch(t *testing.T) {
	client := dial()
	defer client.Close()
	session := client.Session()
	defer session.Release()

	// Set bucket properties.
	// Unfortunately these still aren't exposed via PBC, so do it manually with curl.
	if _, err := exec.Command("curl", "-XPUT", "-H", "content-type:application/json", "http://127.0.0.1:8093/riak/animals", "-d", `{"props":{"precommit":[{"mod":"riak_search_kv_hook","fun":"precommit"}]}}`).Output(); err != nil {
		t.Error(err.Error())
	}

	bucket := session.GetBucket("animals")
	animal := &Animal{
		Name: "pig",
	}
	o1 := bucket.Object("a1")
	o1.ContentType([]byte("application/json"))
	if _, err := o1.Store(animal); err != nil {
		t.Error(err.Error())
	}
	animal = &Animal{
		Name: "dog",
	}
	o2 := bucket.Object("a2")
	o2.ContentType([]byte("application/json"))
	if _, err := o2.Store(animal); err != nil {
		t.Error(err.Error())
	}

	var animals []Animal
	query := session.Query()
	if _, err := query.Key("Id").Out(&animals).Search([]byte("animals"), []byte("name:pig OR name:dog")); err != nil {
		t.Error(err.Error())
	}

	if len(animals) != 2 {
		t.Error("expected 2 documents")
	}

	if animals[0].Id != "a1" && animals[0].Id != "a2" {
		t.Error("unexpected key")
	}

	if _, err := o1.Delete(); err != nil {
		t.Error(err.Error())
	}
	if _, err := o2.Delete(); err != nil {
		t.Error(err.Error())
	}
}
