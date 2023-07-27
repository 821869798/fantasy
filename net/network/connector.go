package network

type Connector struct {
	session ISession
	addr    string
	network INetwork
}

func NewConnector(addr string, network INetwork) *Connector {
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

func (c *Connector) Session() ISession {
	return c.session
}
