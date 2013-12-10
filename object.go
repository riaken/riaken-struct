package riaken_struct

import (
	"errors"
)

import (
	core "github.com/riaken/riaken-core"
	"github.com/riaken/riaken-core/rpb"
)

var ERR_NO_CONTENT error = errors.New("no content")

type Object struct {
	bucket     *Bucket
	coreObject *core.Object
}

// CoreObject returns the underlying riaken-core Object.
func (o *Object) CoreObject() *core.Object {
	return o.coreObject
}

func (o *Object) Do(opts interface{}) *Object {
	o.coreObject.Do(opts)
	return o
}

func (o *Object) ContentType(ct []byte) {
	o.coreObject.ContentType(ct)
}

func (o *Object) Fetch(out interface{}) (*rpb.RpbGetResp, error) {
	data, err := o.coreObject.Fetch()
	if err != nil {
		return nil, err
	}
	if len(data.GetContent()) == 0 {
		return data, ERR_NO_CONTENT
	}
	return data, o.bucket.session.marshaller.Unmarshal(data.GetContent()[0].GetValue(), out)
}

func (o *Object) Store(in interface{}) (*rpb.RpbPutResp, error) {
	content, err := o.bucket.session.marshaller.Marshal(in)
	if err != nil {
		return nil, err
	}
	opts := &rpb.RpbPutReq{
		Content: content,
	}
	return o.coreObject.Do(opts).Store(content.GetValue())
}

func (o *Object) Delete() (bool, error) {
	return o.coreObject.Delete()
}
