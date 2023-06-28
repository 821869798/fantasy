package ws

import (
	"encoding/binary"
	"github.com/821869798/fantasy/net/api"
	"github.com/821869798/fantasy/net/base"
	"github.com/gookit/slog"
	"github.com/gorilla/websocket"
)

type WsConnector struct {
	session api.ISession
	addr    string

	wsSessionAdapter *wsSessionAdapter
}

func NewWsConnector(addr string, handle api.IMsgHandle, codec api.IMsgCodec, opt *WsStartOpt) *WsConnector {
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

	s := base.NewSession(1, conn, c.wsSessionAdapter)
	c.session = s
	s.Start()

}

func (c *WsConnector) Session() api.ISession {
	return c.session
}
