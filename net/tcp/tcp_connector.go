package tcp

import (
	"encoding/binary"
	"github.com/821869798/fantasy/net/api"
	log "github.com/FishGoddess/logit"
	"net"
)

type TcpConnector struct {
	session     api.Session
	addr        string
	transmitter api.MsgTransmitter
	handle      api.MsgHandle
	opt         *TcpStartOpt
}

func NewTcpConnector(addr string, handle api.MsgHandle, transmitter api.MsgTransmitter, opt *TcpStartOpt) *TcpConnector {
	c := &TcpConnector{
		addr:        addr,
		transmitter: transmitter,
		handle:      handle,
		opt:         opt,
	}

	if c.opt == nil {
		// Create Default
		c.opt = NewTcpStartOpt()
	}
	if c.transmitter == nil {
		c.transmitter = NewTcpTransmitter(binary.LittleEndian)
	}

	return c
}

func (c *TcpConnector) Start() {
	go c.run()
}

func (c *TcpConnector) run() {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		log.Error("TcpConnector connect error %v", err)
		return
	}

	s := newTcpSession(1, conn, c.opt.SendChanSize, c.transmitter, c.handle)
	c.session = s
	s.Start()

}

func (c *TcpConnector) Session() api.Session {
	return c.session
}
