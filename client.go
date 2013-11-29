package riaken_struct

import (
	"time"
)

import (
	core "github.com/riaken/riaken-core"
)

type Client struct {
	coreClient *core.Client
	marshaller *StructMarshal
}

func NewClient(addrs []string, max int, timeout time.Duration, sm *StructMarshal) *Client {
	c := new(Client)
	c.coreClient = core.NewClient(addrs, max, timeout)
	c.marshaller = sm
	return c
}

func (c *Client) Dial() {
	c.coreClient.Dial()
}

func (c *Client) Session() (*Session, error) {
	s := new(Session)
	cs, err := c.coreClient.Session()
	if err != nil {
		return nil, err
	}
	s.coreSession = cs
	s.marshaller = c.marshaller
	return s, nil
}

func (c *Client) Close() {
	c.coreClient.Close()
}

// CoreClient fetches the underlying riaken-core Client.
func (c *Client) CoreClient() *core.Client {
	return c.coreClient
}
