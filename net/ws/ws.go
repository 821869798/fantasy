package ws

import (
	"github.com/821869798/fantasy/net/api"
	"github.com/gorilla/websocket"
	"net"
)

type wsSessionAdapter struct {
	codec  api.IMsgCodec
	handle api.IMsgHandle
	opt    *WsStartOpt
}

func (a *wsSessionAdapter) Name() string {
	return "WebSocketSession"
}

func (a *wsSessionAdapter) SendChanSize() uint32 {
	return uint32(a.opt.WriteBufferSize)
}

func (a *wsSessionAdapter) RemoteAddr(rawConn interface{}) net.Addr {
	conn, _ := rawConn.(*websocket.Conn)
	return conn.RemoteAddr()
}

func (a *wsSessionAdapter) CloseConn(rawConn interface{}) error {
	conn, _ := rawConn.(*websocket.Conn)
	return conn.Close()
}

func (a *wsSessionAdapter) Handle() api.IMsgHandle {
	return a.handle
}
func (a *wsSessionAdapter) SendMsg(s api.ISession, msg interface{}) error {
	return a.codec.OnSendMsg(s, msg)
}

func (a *wsSessionAdapter) RecvMsg(s api.ISession) (interface{}, error) {
	msg, err := a.codec.OnRecvMsg(s)
	return msg, err
}
