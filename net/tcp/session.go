package tcp

import "net"

type TcpSession struct {
	conn      net.Conn
}

func (s *TcpSession) sendLoop(){
	
}

func (s *TcpSession) recvLoop(){

}