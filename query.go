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
	out       interface{}
}

func (q *Query) reset() {
	q.out = nil
}

func (q *Query) Do(opts interface{}) *Query {
	q.coreQuery.Do(opts)
	return q
}

func (q *Query) Out(out interface{}) *Query {
	q.out = out
	return q
}

func (q *Query) MapReduce(req, ct []byte) (*rpb.RpbMapRedResp, error) {
	return q.coreQuery.MapReduce(req, ct)
}

// SecondaryIndexes is based loosely on the logic in http://godoc.org/labix.org/v2/mgo#Query.All
//
// This must be called in combination with Out().
func (q *Query) SecondaryIndexes(bucket, index, key, start, end []byte, maxResults uint32, continuation []byte) (*rpb.RpbIndexResp, error) {
	defer q.reset()
	if q.out == nil {
		return nil, errors.New("out must be set via Out()")
	}
	rv := reflect.ValueOf(q.out)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return nil, errors.New("out must be of type slice")
	}
	data, err := q.coreQuery.SecondaryIndexes(bucket, index, key, start, end, maxResults, continuation)
	if err != nil {
		return nil, err
	}
	if len(data.GetKeys()) == 0 {
		return data, nil
	}
	i := 0
	sv := rv.Elem()
	sv = sv.Slice(0, sv.Cap())
	et := sv.Type().Elem()
	b := q.session.GetBucket(string(bucket))
	for _, k := range data.GetKeys() {
		e := reflect.New(et)
		object := b.Object(string(k))
		if _, err := object.Fetch(e.Interface()); err != nil {
			return nil, err
		}
		sv = reflect.Append(sv, e.Elem())
		sv = sv.Slice(0, sv.Cap())
		i++
	}
	rv.Elem().Set(sv.Slice(0, i))
	return data, nil
}

// Search is based loosely on the logic in http://godoc.org/labix.org/v2/mgo#Query.All
// This must be called in combination with Out().
//
// WARNING: Searches that result in a lot of results can potentially run the application out of memory.
func (q *Query) Search(index, query []byte) (*rpb.RpbSearchQueryResp, error) {
	defer q.reset()
	if q.out == nil {
		return nil, errors.New("out must be set via Out()")
	}
	rv := reflect.ValueOf(q.out)
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
