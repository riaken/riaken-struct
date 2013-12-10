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
	key       string
}

func (q *Query) reset() {
	q.out = nil
	q.key = ""
}

// CoreQuery returns the underlying riaken-core Query.
func (q *Query) CoreQuery() *core.Query {
	return q.coreQuery
}

func (q *Query) Do(opts interface{}) *Query {
	q.coreQuery.Do(opts)
	return q
}

func (q *Query) Out(out interface{}) *Query {
	q.out = out
	return q
}

// Key sets the special case struct member to write the key value out to in searches.
func (q *Query) Key(key string) *Query {
	q.key = key
	return q
}

func (q *Query) MapReduce(req, ct []byte) (*rpb.RpbMapRedResp, error) {
	return q.coreQuery.MapReduce(req, ct)
}

// SecondaryIndexes is based loosely on the logic in http://godoc.org/labix.org/v2/mgo#Query.All
//
// Chain call this method with Key() to describe the struct member to write the key value out to,
// and Out() to pass the []struct to output the queried values to.
//
//  type Results struct {
//    Id string
//    Value string
//  }
//
//  var foo []results
//  query.Key("Id").Out(&results).SecondaryIndexes(...)
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
		if q.key != "" {
			ek := reflect.ValueOf(e.Interface()).Elem()
			fk := ek.FieldByName(q.key)
			if fk.Kind() == reflect.String {
				fk.SetString(string(k))
			}
			if fk.Kind() == reflect.Slice {
				fk.SetBytes(k)
			}
		}
		sv = reflect.Append(sv, e.Elem())
		sv = sv.Slice(0, sv.Cap())
		i++
	}
	rv.Elem().Set(sv.Slice(0, i))
	return data, nil
}

// Search is based loosely on the logic in http://godoc.org/labix.org/v2/mgo#Query.All
//
// Chain call this method with Key() to describe the struct member to write the key value out to,
// and Out() to pass the []struct to output the queried values to.
//
// WARNING: Searches that result in a lot of results can potentially run the application out of memory.
// If this occurs fall back to fetching the keys with the riaken-core search and manually fetching the data.
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

	// Since a key can be returned multiple times, check if it's already been stored.
	found := map[string]bool{}
	check := func(key string) bool {
		_, ok := found[key]
		if !ok {
			found[key] = true
		}
		return ok
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
				// Check if key has been stored.
				if check(string(v.GetValue())) {
					continue
				}
				e := reflect.New(et)
				object := bucket.Object(string(v.GetValue()))
				if _, err := object.Fetch(e.Interface()); err != nil {
					return nil, err
				}
				if q.key != "" {
					ek := reflect.ValueOf(e.Interface()).Elem()
					fk := ek.FieldByName(q.key)
					if fk.Kind() == reflect.String {
						fk.SetString(string(v.GetValue()))
					}
					if fk.Kind() == reflect.Slice {
						fk.SetBytes(v.GetValue())
					}
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
