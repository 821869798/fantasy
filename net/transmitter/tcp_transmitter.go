package transmitter

import (
	"encoding/binary"
	"github.com/821869798/fantasy/net/api"
)

type tcpTransmitter struct {
	order binary.ByteOrder
}

func NewTcpTransmitter(order binary.ByteOrder) api.MsgTransmitter{
	t := &tcpTransmitter{
		order: order,
	}
	return t
}

func (t *tcpTransmitter)OnSendMsg(s api.Session,msg interface{}) error{
	return nil
}
func (t *tcpTransmitter)OnRecvMsg(s api.Session) (interface{},error){
	return nil,nil
}