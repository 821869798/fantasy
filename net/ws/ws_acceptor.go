package ws

import (
	"context"
	"encoding/binary"
	"github.com/821869798/fantasy/net/network"
	"github.com/gookit/slog"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type WsAcceptor struct {
	addr string

	wsUpgrader websocket.Upgrader

	wsSessionAdapter *wsSessionAdapter

	sessionIDGen uint64
	sessMap      sync.Map

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewWSAcceptor(addr string, handle network.IMsgHandle, codec network.IMsgCodec, opt *WsStartOpt) *WsAcceptor {

	if opt == nil {
		// Create Default
		opt = NewWsStartOpt()
	}
	if codec == nil {
		codec = NewWsMsgCodec(binary.BigEndian)
	}

	wsSessionAdapter := &wsSessionAdapter{
		handle: handle,
		codec:  codec,
		opt:    opt,
	}

	a := &WsAcceptor{
		addr:             addr,
		wsSessionAdapter: wsSessionAdapter,
		wsUpgrader: websocket.Upgrader{
			ReadBufferSize:  opt.ReadBufferSize,
			WriteBufferSize: opt.WriteBufferSize,
		},
	}

	a.ctx, a.ctxCancel = context.WithCancel(context.Background())

	return a
}

func (t *WsAcceptor) Start() {

	if !t.init() {
		return
	}

	go t.run()
}

func (t *WsAcceptor) init() bool {
	// Listen
	http.HandleFunc(t.wsSessionAdapter.opt.HandlePath, func(w http.ResponseWriter, r *http.Request) {
		conn, err := t.wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			// handle error
			return
		}
		t.sessionIDGen++
		sid := t.sessionIDGen
		go t.handleSession(sid, conn)
	})

	return true
}

func (t *WsAcceptor) run() {
	if err := http.ListenAndServe(t.addr, nil); err != nil {
		if err != nil {
			slog.Errorf("Http Listen error:%v", err)
			return
		}
	}
}

func (t *WsAcceptor) handleSession(sid uint64, conn *websocket.Conn) {

	//handle session
	s := network.NewSession(sid, conn, t.wsSessionAdapter)

	t.sessMap.Store(sid, conn)
	s.Start()
	t.sessMap.Delete(sid)
}

func (t *WsAcceptor) GetSession(sid uint64) network.ISession {
	v, ok := t.sessMap.Load(sid)
	if !ok {
		return nil
	}
	s, ok := v.(network.ISession)
	if !ok {
		return nil
	}
	return s
}
