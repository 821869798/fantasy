package tcp

import (
	"encoding/binary"
	"github.com/821869798/fantasy/net/api"
	"github.com/821869798/fantasy/net/base"
	"github.com/821869798/fantasy/net/packet"
	"github.com/gookit/slog"
	"golang.org/x/net/netutil"
	"net"
)

func NewTcpAcceptor(addr string, handle api.IMsgHandle, codec api.IMsgCodec, opt *TcpStartOpt) *base.Acceptor {
	tcpNetwork := newTcpNetwork(handle, codec, opt)
	return base.NewAcceptor(addr, tcpNetwork)
}

func NewTcpConnector(addr string, handle api.IMsgHandle, codec api.IMsgCodec, opt *TcpStartOpt) *base.Connector {
	tcpNetwork := newTcpNetwork(handle, codec, opt)
	return base.NewConnector(addr, tcpNetwork)
}

type tcpNetwork struct {
	opt     *TcpStartOpt
	adapter *tcpSessionAdapter
}

func newTcpNetwork(handle api.IMsgHandle, codec api.IMsgCodec, opt *TcpStartOpt) *tcpNetwork {

	if opt == nil {
		// Create Default
		opt = NewTcpStartOpt()
	}
	if codec == nil {
		codec = packet.NewLTVMsgCodec(binary.BigEndian)
	}
	n := &tcpNetwork{
		opt: opt,
		adapter: &tcpSessionAdapter{
			opt:    opt,
			codec:  codec,
			handle: handle,
		},
	}
	return n
}

func (n *tcpNetwork) CreateListener(addr string) (net.Listener, bool) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Errorf("TcpAcceptor Listen error %v", err)
		return nil, false
	}
	//limit connect count
	listener := netutil.LimitListener(ln, n.opt.MaxConns)
	return listener, true
}
func (n *tcpNetwork) Dial(addr string) (net.Conn, bool) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		slog.Errorf("TcpConnector connect error %v", err)
		return nil, false
	}
	return conn, true
}

func (n *tcpNetwork) SessionAdapter() api.ISessionAdapter {
	return n.adapter
}

type tcpSessionAdapter struct {
	codec  api.IMsgCodec
	handle api.IMsgHandle
	opt    *TcpStartOpt
}

func (a *tcpSessionAdapter) Name() string {
	return "TcpSession"
}

func (a *tcpSessionAdapter) SendChanSize() uint32 {
	return a.opt.SendChanSize
}
func (a *tcpSessionAdapter) Handle() api.IMsgHandle {
	return a.handle
}
func (a *tcpSessionAdapter) SendMsg(s api.ISession, msg interface{}) error {
	return a.codec.OnSendMsg(s, msg)
}

func (a *tcpSessionAdapter) RecvMsg(s api.ISession) (interface{}, error) {
	msg, err := a.codec.OnRecvMsg(s)
	return msg, err
}
