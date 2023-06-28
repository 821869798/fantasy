package base

import (
	"github.com/821869798/fantasy/net/api"
)

type Connector struct {
	session api.ISession
	addr    string
	network api.INetwork
}

func NewConnector(addr string, network api.INetwork) *Connector {
	c := &Connector{
		addr:    addr,
		network: network,
	}

	return c
}

func (c *Connector) Start() {
	go c.run()
}

func (c *Connector) run() {
	conn, ok := c.network.Dial(c.addr)
	if !ok {
		return
	}

	s := NewSession(1, conn, c.network.SessionAdapter())
	c.session = s
	s.Start()

}

func (c *Connector) Session() api.ISession {
	return c.session
}
