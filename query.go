package riaken_struct

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

func (q *Query) Search(index, query []byte) (*rpb.RpbSearchQueryResp, error) {
	return q.coreQuery.Search(index, query)
}
