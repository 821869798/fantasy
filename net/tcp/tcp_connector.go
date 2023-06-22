package tcp

import (
	"encoding/binary"
	"github.com/821869798/fantasy/net/api"
	"github.com/gookit/slog"
	"net"
)

type TcpConnector struct {
	session api.Session
	addr    string
	codec   api.MsgCodec
	handle  api.MsgHandle
	opt     *TcpStartOpt
}

func NewTcpConnector(addr string, handle api.MsgHandle, codec api.MsgCodec, opt *TcpStartOpt) *TcpConnector {
	c := &TcpConnector{
		addr:   addr,
		codec:  codec,
		handle: handle,
		opt:    opt,
	}

	if c.opt == nil {
		// Create Default
		c.opt = NewTcpStartOpt()
	}
	if c.codec == nil {
		c.codec = NewTcpMsgCodec(binary.BigEndian)
	}

	return c
}

func (c *TcpConnector) Start() {
	go c.run()
}

func (c *TcpConnector) run() {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		slog.Errorf("TcpConnector connect error %v", err)
		return
	}

	s := newTcpSession(1, conn, c.opt.SendChanSize, c.codec, c.handle)
	c.session = s
	s.Start()

}

func (c *TcpConnector) Session() api.Session {
	return c.session
}
