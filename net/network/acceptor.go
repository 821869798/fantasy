package network

import (
	"context"
	"github.com/gookit/slog"
	"net"
	"sync"
)

type Acceptor struct {
	addr string

	sessionIDGen uint64
	sessMap      sync.Map

	ctx       context.Context
	ctxCancel context.CancelFunc

	listener net.Listener

	network INetwork
}

func NewAcceptor(addr string, network INetwork) *Acceptor {
	a := &Acceptor{
		addr:    addr,
		network: network,
	}

	a.ctx, a.ctxCancel = context.WithCancel(context.Background())

	return a
}

func (t *Acceptor) Start() {

	if !t.init() {
		return
	}

	go t.run()
}

func (t *Acceptor) init() bool {

	listener, ok := t.network.CreateListener(t.addr)
	if ok {
		t.listener = listener
	}
	return ok
}

func (t *Acceptor) run() {

	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			{
				conn, err := t.listener.Accept()
				if err != nil {
					slog.Errorf("Acceptor accept connection error:%v", err)
					continue
				}
				t.sessionIDGen++
				sid := t.sessionIDGen
				go t.handleSession(sid, conn)
			}
		}
	}

}

func (t *Acceptor) handleSession(sid uint64, conn net.Conn) {
	s := NewSession(sid, conn, t.network.SessionAdapter())
	t.sessMap.Store(sid, conn)
	s.Start()
	t.sessMap.Delete(sid)
}

func (t *Acceptor) GetSession(sid uint64) ISession {
	v, ok := t.sessMap.Load(sid)
	if !ok {
		return nil
	}
	s, ok := v.(ISession)
	if !ok {
		return nil
	}
	return s
}
