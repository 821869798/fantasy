package kcp

import (
	"encoding/binary"
	"github.com/821869798/fantasy/net/network"
	"github.com/821869798/fantasy/net/packet"
	"github.com/gookit/slog"
	"github.com/xtaci/kcp-go/v5"
	"golang.org/x/net/netutil"
	"net"
	"time"
)

func NewKcpAcceptor(addr string, handle network.IMsgHandle, codec network.IMsgCodec, opt *KcpStartOpt) *network.Acceptor {
	kcpNetwork := newKcpNetwork(handle, codec, opt)
	return network.NewAcceptor(addr, kcpNetwork)
}

func NewKcpConnector(addr string, handle network.IMsgHandle, codec network.IMsgCodec, opt *KcpStartOpt) *network.Connector {
	kcpNetwork := newKcpNetwork(handle, codec, opt)
	return network.NewConnector(addr, kcpNetwork)
}

type kcpNetwork struct {
	opt     *KcpStartOpt
	adapter *kcpSessionAdapter
}

func newKcpNetwork(handle network.IMsgHandle, codec network.IMsgCodec, opt *KcpStartOpt) *kcpNetwork {

	if opt == nil {
		// Create Default
		opt = NewKcpStartOpt()
	}
	if codec == nil {
		codec = packet.NewLTVMsgCodec(binary.BigEndian)
	}
	n := &kcpNetwork{
		opt: opt,
		adapter: &kcpSessionAdapter{
			opt:    opt,
			codec:  codec,
			handle: handle,
		},
	}
	return n
}

func (n *kcpNetwork) CreateListener(addr string) (net.Listener, bool) {
	ln, err := kcp.Listen(addr)
	if err != nil {
		slog.Errorf("KcpAcceptor Listen error %v", err)
		return nil, false
	}
	//limit connect count
	listener := netutil.LimitListener(ln, n.opt.MaxConns)
	return listener, true
}
func (n *kcpNetwork) Dial(addr string) (net.Conn, bool) {
	conn, err := kcp.Dial(addr)
	if err != nil {
		slog.Errorf("KcpConnector connect error %v", err)
		return nil, false
	}
	return conn, true
}

func (n *kcpNetwork) SessionAdapter() network.ISessionAdapter {
	return n.adapter
}

type kcpSessionAdapter struct {
	codec  network.IMsgCodec
	handle network.IMsgHandle
	opt    *KcpStartOpt
}

func (a *kcpSessionAdapter) Name() string {
	return "KcpSession"
}

func (a *kcpSessionAdapter) SendChanSize() uint32 {
	return a.opt.SendChanSize
}

func (a *kcpSessionAdapter) RemoteAddr(rawConn interface{}) net.Addr {
	conn, _ := rawConn.(net.Conn)
	return conn.RemoteAddr()
}

func (a *kcpSessionAdapter) CloseConn(rawConn interface{}) error {
	conn, _ := rawConn.(net.Conn)
	return conn.Close()
}

func (a *kcpSessionAdapter) Handle() network.IMsgHandle {
	return a.handle
}
func (a *kcpSessionAdapter) SendMsg(s network.ISession, msg interface{}) error {
	return a.codec.OnSendMsg(s, msg)
}

func (a *kcpSessionAdapter) RecvMsg(s network.ISession) (interface{}, error) {
	conn, _ := s.Raw().(net.Conn)
	err := conn.SetReadDeadline(time.Now().Add(a.opt.Timeout))
	if err != nil {
		return nil, err
	}
	msg, err := a.codec.OnRecvMsg(s)
	return msg, err
}
