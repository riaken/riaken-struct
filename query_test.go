package riaken_struct

import (
	"os/exec"
	"testing"
)

type Animal struct {
	Name string `json:"name"`
}

func TestQuery(t *testing.T) {
	client := dial()
	defer client.Close()
	session, err := client.Session()
	if err != nil {
		t.Error(err.Error())
	}
	defer session.Close()

	session.Query()
}

func TestSearch(t *testing.T) {
	client := dial()
	defer client.Close()
	session, err := client.Session()
	if err != nil {
		t.Error(err.Error())
	}
	defer session.Close()

	// Set bucket properties.
	// Unfortunately these still aren't exposed via PBC, so do it manually with curl.
	_, err = exec.Command("curl", "-XPUT", "-H", "content-type:application/json", "http://127.0.0.1:8093/riak/animals", "-d", `{"props":{"precommit":[{"mod":"riak_search_kv_hook","fun":"precommit"}]}}`).Output()
	if err != nil {
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
	if _, err := query.Search([]byte("animals"), []byte("name:pig OR name:dog"), &animals); err != nil {
		t.Error(err.Error())
	}

	if len(animals) != 2 {
		t.Error("expected 2 documents")
	}

	if _, err := o1.Delete(); err != nil {
		t.Error(err.Error())
	}
	if _, err := o2.Delete(); err != nil {
		t.Error(err.Error())
	}
}
