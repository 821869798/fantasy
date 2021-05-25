package tcp

import (
	"github.com/821869798/fantasy/net/api"
	"go.uber.org/atomic"
	"net"
)

type TcpAcceptor struct {
	gsid atomic.Uint64
	addr string

	listener net.Listener
}

func NewTcpAcceptor(addr string, transmitter api.MsgTransmitter) {

}

func (t *TcpAcceptor) Start() {
}
