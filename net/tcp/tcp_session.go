package tcp

import (
	"context"
	"errors"
	"github.com/821869798/fantasy/net/api"
	"go.uber.org/atomic"
	"net"
	"sync"
)

var SessionClosedError = errors.New("Session Closed")
var SessionBlockedError = errors.New("Session Blocked")

type tcpSession struct {
	sid  uint64
	conn net.Conn

	sendChan chan interface{}

	transmitter api.MsgTransmitter

	//退出通知
	ctx       context.Context
	ctxCancel context.CancelFunc

	// 退出同步器
	exitSync sync.WaitGroup

	isClose atomic.Bool

	endNotify func()
}

func newTcpSession(sid uint64, conn net.Conn, sendChanSize uint32, transmitter api.MsgTransmitter, endNotify func()) *tcpSession {
	s := &tcpSession{
		sid:         sid,
		conn:        conn,
		sendChan:    make(chan interface{}, sendChanSize),
		transmitter: transmitter,
		endNotify:   endNotify,
	}

	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	return s
}

func (s *tcpSession) Start() {

	s.exitSync.Add(2)

	go func() {
		s.exitSync.Wait()

		if s.endNotify != nil {
			s.endNotify()
		}

	}()

	go s.recvLoop()
	go s.sendLoop()

}

func (s *tcpSession) Raw() interface{} {
	return s.conn
}

func (s *tcpSession) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *tcpSession) Sid() uint64 {
	return s.sid
}

func (s *tcpSession) Send(msg interface{}) error {
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

func (s *tcpSession) sendLoop() {

	defer func() {
		s.Close()
		s.exitSync.Done()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case msg := <-s.sendChan:
			err := s.transmitter.OnSendMsg(s, msg)
			if err != nil {
				return
			}
		}
	}
}

func (s *tcpSession) recvLoop() {

	defer func() {
		s.Close()
		s.exitSync.Done()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		}
	}

}

func (s *tcpSession) IsClose() bool {
	return s.isClose.Load()
}

func (s *tcpSession) Close() {
	if s.IsClose() {
		return
	}

	s.ctxCancel()
	_ = s.conn.Close()
}
