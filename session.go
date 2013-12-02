package riaken_struct

import (
	core "github.com/riaken/riaken-core"
	"github.com/riaken/riaken-core/rpb"
)

type Session struct {
	client      *Client
	coreSession *core.Session
	marshaller  *StructMarshal
}

// CoreSession returns the underlying riaken-core Session.
func (s *Session) CoreSession() *core.Session {
	return s.coreSession
}

func (s *Session) Release() {
	s.coreSession.Release()
}

func (s *Session) GetBucket(name string) *Bucket {
	b := new(Bucket)
	b.session = s
	b.coreBucket = s.coreSession.GetBucket(name)
	return b
}

func (s *Session) Query() *Query {
	q := new(Query)
	q.session = s
	q.coreQuery = s.coreSession.Query()
	return q
}

func (s *Session) ListBuckets() ([]*core.Bucket, error) {
	return s.coreSession.ListBuckets()
}

func (s *Session) Ping() bool {
	return s.coreSession.Ping()
}

func (s *Session) GetClientId() (*rpb.RpbGetClientIdResp, error) {
	return s.coreSession.GetClientId()
}

func (s *Session) SetClientId(id []byte) (bool, error) {
	return s.coreSession.SetClientId(id)
}

func (s *Session) ServerInfo() (*rpb.RpbGetServerInfoResp, error) {
	return s.coreSession.ServerInfo()
}
