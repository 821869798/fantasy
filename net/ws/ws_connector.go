package ws

import (
	"encoding/binary"
	"github.com/821869798/fantasy/net/network"
	"github.com/gookit/slog"
	"github.com/gorilla/websocket"
)

type WsConnector struct {
	session network.ISession
	addr    string

	wsSessionAdapter *wsSessionAdapter
}

func NewWsConnector(addr string, handle network.IMsgHandle, codec network.IMsgCodec, opt *WsStartOpt) *WsConnector {
	c := &WsConnector{
		addr: addr,
	}

	if opt == nil {
		// Create Default
		opt = NewWsStartOpt()
	}
	if codec == nil {
		codec = NewWsMsgCodec(binary.BigEndian)
	}

	c.wsSessionAdapter = &wsSessionAdapter{
		handle: handle,
		codec:  codec,
		opt:    opt,
	}

	return c
}

func (c *WsConnector) Start() {
	go c.run()
}

func (c *WsConnector) run() {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(c.addr, nil)
	if err != nil {
		slog.Errorf("WebSocketConnector connect error %v", err)
		return
	}

	s := network.NewSession(1, conn, c.wsSessionAdapter)
	c.session = s
	s.Start()

}

func (c *WsConnector) Session() network.ISession {
	return c.session
}
