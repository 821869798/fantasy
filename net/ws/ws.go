package ws

import (
	"github.com/821869798/fantasy/net/network"
	"github.com/gorilla/websocket"
	"net"
)

type wsSessionAdapter struct {
	codec  network.IMsgCodec
	handle network.IMsgHandle
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

func (a *wsSessionAdapter) Handle() network.IMsgHandle {
	return a.handle
}
func (a *wsSessionAdapter) SendMsg(s network.ISession, msg interface{}) error {
	return a.codec.OnSendMsg(s, msg)
}

func (a *wsSessionAdapter) RecvMsg(s network.ISession) (interface{}, error) {
	msg, err := a.codec.OnRecvMsg(s)
	return msg, err
}
