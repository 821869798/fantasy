package tcp

import (
	"context"
	"encoding/binary"
	"github.com/821869798/fantasy/net/api"
	"github.com/gookit/slog"
	"go.uber.org/atomic"
	"golang.org/x/net/netutil"
	"net"
	"sync"
)

type TcpAcceptor struct {
	gsid   atomic.Uint64
	addr   string
	codec  api.MsgCodec
	handle api.MsgHandle
	opt    *TcpStartOpt

	sessionIDGen uint64
	sessMap      sync.Map

	ctx       context.Context
	ctxCancel context.CancelFunc

	listener net.Listener
}

func NewTcpAcceptor(addr string, handle api.MsgHandle, codec api.MsgCodec, opt *TcpStartOpt) *TcpAcceptor {
	a := &TcpAcceptor{
		addr:   addr,
		codec:  codec,
		handle: handle,
		opt:    opt,
	}

	a.ctx, a.ctxCancel = context.WithCancel(context.Background())
	if a.opt == nil {
		// Create Default
		a.opt = NewTcpStartOpt()
	}
	if a.codec == nil {
		a.codec = NewTcpMsgCodec(binary.BigEndian)
	}

	return a
}

func (t *TcpAcceptor) Start() {

	if !t.init() {
		return
	}

	go t.run()
}

func (t *TcpAcceptor) init() bool {

	ln, err := net.Listen("tcp", t.addr)
	if err != nil {
		slog.Errorf("TcpAcceptor Listen error %v", err)
		return false
	}
	//limit connect count
	t.listener = netutil.LimitListener(ln, t.opt.MaxConns)
	return true
}

func (t *TcpAcceptor) run() {

	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			{
				conn, err := t.listener.Accept()
				if err != nil {
					slog.Errorf("TcpAcceptor accept connection error:%v", err)
					continue
				}
				t.sessionIDGen++
				sid := t.sessionIDGen
				go t.handleSession(sid, conn)
			}
		}
	}

}

func (t *TcpAcceptor) handleSession(sid uint64, conn net.Conn) {
	slog.Debugf("TcpAcceptor handle new session sid %v addr %v", sid, conn.RemoteAddr())

	s := newTcpSession(sid, conn, t.opt.SendChanSize, t.codec, t.handle)
	t.sessMap.Store(sid, conn)
	s.Start()
	t.sessMap.Delete(sid)
}

func (t *TcpAcceptor) GetSession(sid uint64) api.Session {
	v, ok := t.sessMap.Load(sid)
	if !ok {
		return nil
	}
	s, ok := v.(*tcpSession)
	if !ok {
		return nil
	}
	return s
}
