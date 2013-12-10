package riaken_struct

import (
	core "github.com/riaken/riaken-core"
	"github.com/riaken/riaken-core/rpb"
)

type Bucket struct {
	session    *Session
	coreBucket *core.Bucket
}

// CoreBucket returns the underlying riaken-core Bucket.
func (b *Bucket) CoreBucket() *core.Bucket {
	return b.coreBucket
}

func (b *Bucket) Object(key string) *Object {
	o := new(Object)
	o.bucket = b
	o.coreObject = b.coreBucket.Object(key)
	return o
}

func (b *Bucket) ListKeys() (*rpb.RpbListKeysResp, error) {
	return b.coreBucket.ListKeys()
}

func (b *Bucket) GetBucketProps() (*rpb.RpbGetBucketResp, error) {
	return b.coreBucket.GetBucketProps()
}

func (b *Bucket) SetBucketProps(props *rpb.RpbBucketProps) (bool, error) {
	return b.coreBucket.SetBucketProps(props)
}
