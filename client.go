package riaken_struct

import (
	core "github.com/riaken/riaken-core"
)

type Client struct {
	coreClient *core.Client
	marshaller *StructMarshal
}

func NewClient(addrs []string, max int, sm *StructMarshal) *Client {
	c := new(Client)
	c.coreClient = core.NewClient(addrs, max)
	c.marshaller = sm
	return c
}

// CoreClient fetches the underlying riaken-core Client.
func (c *Client) CoreClient() *core.Client {
	return c.coreClient
}

func (c *Client) Dial() error {
	return c.coreClient.Dial()
}

func (c *Client) Debug(debug bool) {
	c.coreClient.Debug(debug)
}

func (c *Client) Session() *Session {
	s := new(Session)
	s.coreSession = c.coreClient.Session()
	s.marshaller = c.marshaller
	return s
}

func (c *Client) Close() {
	c.coreClient.Close()
}
