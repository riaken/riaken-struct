## Riaken Struct

High level struct wrapper for riaken-core.

It provides the ability to pass structs in and out of Riak as other formats such as JSON, XML, etc.

## Install

    go get gopkg.in/riaken/riaken-struct.v1

## Documentation

http://godoc.org/gopkg.in/riaken/riaken-struct.v1

The following methods being wrapped for struct usage are:

* Object->Fetch
* Object->Store
* Query->SecondaryIndexes
* Query->Search

The rest of the underlying riaken-core API is exposed as normal.

## Usage

### Client

Not much different from riaken-core, with the exception of passing a marshaller.

	package main

	import "log"
	import "github.com/riaken/riaken-struct"

	func main() {
		// Marshaller
		marshaller := NewStructMarshal("json", JsonMarshaller, JsonUnmarshaller)
		// Riak cluster addresses
		addrs := []string{"127.0.0.1:8083", "127.0.0.1:8084", "127.0.0.1:8085", "127.0.0.1:8086", "127.0.0.1:8087"}
		// Create a client, passing the addresses, and number of connections to maintain per cluster node
		client := riaken_struct.NewClient(addrs, 1, marshaller)
		// Dial the servers
		if err := client.Dial(); err != nil {
			log.Fatal(err.Error()) // all nodes are down
		}
		// Gracefully close all the connections when finished
		defer client.Close()

		// Grab a session to interact with the cluster
		session := client.Session()
		// Release the session
		defer session.Release()
	}

The marshaller allows structs with tag data to be parsed to and from their target format, in this case JSON.  If a format other than JSON is desired consider using `struct_marshall.go` as a guideline.

Note that all methods can utilize the same `Do()` construct as riaken-core.

### Object

Standard struct data.

	type User struct {
		Name string `json:"name"`
	}

#### Store

	user := &User{
		Name: "User",
	}
	bucket := session.GetBucket("b1")
	object := bucket.Object("o1")
	if _, err := object.Store(user); err != nil {
		log.Error(err.Error())
	}

#### Fetch

	check := User{}
	if _, err := object.Fetch(&check); err != nil {
		log.Error(err.Error())
	}
	log.Print(check.Name) // User

### Query

Since we are querying out fields where we don't know the keys, there are a few things to note.

	type Animal struct {
		Id   string
		Name string `json:"name" riak:"index"`
	}

* `Id` is arbitrarily named, but necessary.  It will contain the key from the server.
* `riak:"index"` allows us to specify this member as a queryable 2i index.

Note that each method below chains with `Key()` and `Out()`.  The first specifies the struct member in which to place the key, the second the struct slice to output results to.

These can potentially run the application out of memory, so use with care.  As opposed to their lower level riaken-core counterparts they actively fetch the object for every key found as well (without duplicates).  If lower level behavior is desired through this driver consider calling `session.Query().CoreQuery()` and manually fetching objects as required.

#### Secondary Indexes

	var check []Animal
	query := session.Query()
	if _, err := query.Key("Id").Out(&check).SecondaryIndexes([]byte("animals"), []byte("name_bin"), []byte("chicken"), nil, nil, 0, nil); err != nil {
		log.Error(err.Error())
	}

#### Search

	var animals []Animal
	query := session.Query()
	if _, err := query.Key("Id").Out(&animals).Search([]byte("animals"), []byte("name:pig OR name:dog")); err != nil {
		log.Error(err.Error())
	}

## Author

Brian Jones - mojobojo@gmail.com - https://twitter.com/mojobojo

## License

http://boj.mit-license.org/
