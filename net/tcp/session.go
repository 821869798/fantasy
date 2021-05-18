package tcp

import (
	"context"
	"go.uber.org/atomic"
	"net"
)

type tcpSession struct {
	sid uint64
	conn      net.Conn

	ctx context.Context
	ctxCancel context.CancelFunc

	isClose    atomic.Bool

}

func newTcpSession(sid uint64,conn net.Conn) *tcpSession{
	s := &tcpSession{
		sid : sid,
		conn : conn,
	}

	s.ctx , s.ctxCancel = context.WithCancel(context.Background())
	return s
}

func (s *tcpSession) Sid()uint64{
	return s.sid
}


func (s *tcpSession) Send(msg interface{}){
	if s.IsClose() {
		return
	}



}

func (s *tcpSession) sendLoop(){
	for {
		select {
			case <- s.ctx.Done():
				return
		}
	}
}

func (s *tcpSession) recvLoop(){
	for{
		select {
			case <- s.ctx.Done():
				return
		}
	}
}

func (s *tcpSession) IsClose() bool{
	return s.isClose.Load()
}

func (s *tcpSession) Close(){
	if s.IsClose() {
		return
	}
}