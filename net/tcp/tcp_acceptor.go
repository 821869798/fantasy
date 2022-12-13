package tcp

import (
	"context"
	"encoding/binary"
	"github.com/821869798/fantasy/net/api"
	log "github.com/FishGoddess/logit"
	"go.uber.org/atomic"
	"golang.org/x/net/netutil"
	"net"
	"sync"
)

type TcpAcceptor struct {
	gsid        atomic.Uint64
	addr        string
	transmitter api.MsgTransmitter
	handle      api.MsgHandle
	opt         *TcpStartOpt

	sessionIDGen uint64
	sessMap      sync.Map

	ctx       context.Context
	ctxCancel context.CancelFunc

	listener net.Listener
}

func NewTcpAcceptor(addr string, handle api.MsgHandle, transmitter api.MsgTransmitter, opt *TcpStartOpt) *TcpAcceptor {
	a := &TcpAcceptor{
		addr:        addr,
		transmitter: transmitter,
		handle:      handle,
		opt:         opt,
	}

	a.ctx, a.ctxCancel = context.WithCancel(context.Background())
	if a.opt == nil {
		// Create Default
		a.opt = NewTcpStartOpt()
	}
	if a.transmitter == nil {
		a.transmitter = NewTcpTransmitter(binary.LittleEndian)
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
		log.Error("TcpAcceptor Listen error %v", err)
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
					log.Error("TcpAcceptor accept connection error ", err)
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
	log.Debug("TcpAcceptor handle new session sid %v addr %v", sid, conn.RemoteAddr())

	s := newTcpSession(sid, conn, t.opt.SendChanSize, t.transmitter, t.handle)
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
