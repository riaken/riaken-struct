## Riaken Struct

High level struct wrapper for riaken-core.

It provides the ability to pass structs in and out of Riak as other formats such as JSON, XML, etc.

### Install

    go get github.com/riaken/riaken-struct

### Documentation

http://godoc.org/github.com/riaken/riaken-struct

The following methods being wrapped for struct usage are:

* Object->Fetch
* Object->Store
* Query->SecondaryIndexes
* Query->Search

The rest of the underlying riaken-core API is exposed as normal.

### Author

Brian Jones - mojobojo@gmail.com

