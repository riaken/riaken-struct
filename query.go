package riaken_struct

import (
	"errors"
	"reflect"
)

import (
	core "github.com/riaken/riaken-core"
	"github.com/riaken/riaken-core/rpb"
)

type Query struct {
	session   *Session
	coreQuery *core.Query
}

func (q *Query) Do(opts interface{}) *Query {
	q.coreQuery.Do(opts)
	return q
}

func (q *Query) MapReduce(req, ct []byte) (*rpb.RpbMapRedResp, error) {
	return q.coreQuery.MapReduce(req, ct)
}

func (q *Query) SecondaryIndexes(bucket, index, key, start, end []byte, maxResults uint32, continuation []byte) (*rpb.RpbIndexResp, error) {
	return q.coreQuery.SecondaryIndexes(bucket, index, key, start, end, maxResults, continuation)
}

// Search is based loosely on the logic in http://godoc.org/labix.org/v2/mgo#Query.All
//
// WARNING: Searches that result in a lot of results can potentially run the application out of memory.
func (q *Query) Search(index, query []byte, out interface{}) (*rpb.RpbSearchQueryResp, error) {
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return nil, errors.New("out must be of type slice")
	}
	data, err := q.coreQuery.Search(index, query)
	if err != nil {
		return nil, err
	}
	if data.GetNumFound() == 0 {
		return data, nil
	}
	i := 0
	sv := rv.Elem()
	sv = sv.Slice(0, sv.Cap())
	et := sv.Type().Elem()
	bucket := q.session.GetBucket(string(index))
	for _, d := range data.GetDocs() {
		for _, v := range d.GetFields() {
			// Riak seems to default the key to "id"
			if string(v.GetKey()) == "id" {
				e := reflect.New(et)
				object := bucket.Object(string(v.GetValue()))
				if _, err := object.Fetch(e.Interface()); err != nil {
					return nil, err
				}
				sv = reflect.Append(sv, e.Elem())
				sv = sv.Slice(0, sv.Cap())
				i++
			}
		}
	}
	rv.Elem().Set(sv.Slice(0, i))
	return data, nil
}
