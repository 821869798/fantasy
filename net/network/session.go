package network

import (
	"context"
	"github.com/gookit/slog"
	"go.uber.org/atomic"
	"net"
	"sync"
)

type Session struct {
	sid      uint64
	conn     interface{}
	sendChan chan interface{}

	adapter ISessionAdapter
	handle  IMsgHandle

	//退出通知
	ctx       context.Context
	ctxCancel context.CancelFunc

	// 退出同步器
	exitSync sync.WaitGroup

	isClose atomic.Bool
}

func NewSession(sid uint64, conn interface{}, adapter ISessionAdapter) *Session {
	s := &Session{
		sid:      sid,
		conn:     conn,
		sendChan: make(chan interface{}, adapter.SendChanSize()),
		adapter:  adapter,
		handle:   adapter.Handle(),
	}

	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	return s
}

func (s *Session) Start() {

	slog.Debugf("%s[%v] created,sid:%v", s.adapter.Name(), s.RemoteAddr(), s.sid)

	s.handle.TriggerEvent(&SessionAdd{Session: s})

	s.exitSync.Add(2)

	go s.recvLoop()
	go s.sendLoop()

	s.exitSync.Wait()

	s.handle.TriggerEvent(&SessionRemove{Session: s})
}

func (s *Session) Raw() interface{} {
	return s.conn
}

func (s *Session) RemoteAddr() net.Addr {
	return s.adapter.RemoteAddr(s.conn)
}

func (s *Session) Sid() uint64 {
	return s.sid
}

func (s *Session) Send(msg interface{}) error {
	if s.IsClose() {
		return SessionClosedError
	}

	select {
	case s.sendChan <- msg:
		return nil
	default:
		s.Close()
		return SessionBlockedError
	}
}

func (s *Session) sendLoop() {

	defer func() {
		s.Close()
		s.exitSync.Done()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case msg := <-s.sendChan:
			err := s.adapter.SendMsg(s, msg)
			if err != nil {
				slog.Warnf("%s[%v] send msg error %v", s.adapter.Name(), s.RemoteAddr(), err)
				return
			}
		}
	}
}

func (s *Session) recvLoop() {

	defer func() {
		s.Close()
		s.exitSync.Done()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		msg, err := s.adapter.RecvMsg(s)

		if err != nil {
			slog.Debugf("%s[%v] receive loop msg error %v", s.adapter.Name(), s.RemoteAddr(), err)
			return
		}

		s.handle.TriggerEvent(&SessionMsg{Session: s, Msg: msg})
	}

}

func (s *Session) IsClose() bool {
	return s.isClose.Load()
}

func (s *Session) Close() {
	if s.IsClose() {
		return
	}

	s.ctxCancel()
	_ = s.adapter.CloseConn(s.conn)

}
