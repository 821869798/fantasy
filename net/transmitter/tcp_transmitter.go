package transmitter

import (
	"encoding/binary"
	"fmt"
	"github.com/821869798/fantasy/net/api"
	"io"
)

const (
	MsgSizeLen = 4
	MsgTypeLen = 4
	MsgHeadLen = MsgSizeLen + MsgTypeLen
	MsgMaxSize = 65536
)

type LTVPacket struct {
	Len uint32
	Type   uint32
	Value   []byte
}

type tcpTransmitter struct {
	order binary.ByteOrder
	headLen []byte
}

func NewTcpTransmitter(order binary.ByteOrder) api.MsgTransmitter{
	t := &tcpTransmitter{
		order: order,
		headLen : make([]byte,MsgSizeLen),
	}
	return t
}

func (t *tcpTransmitter)OnSendMsg(s api.Session,msg interface{}) error{
	return nil
}
func (t *tcpTransmitter)OnRecvMsg(s api.Session) (interface{},error){
	reader, ok := s.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	if _, err := io.ReadFull(reader, t.headLen); err != nil {
		return nil, err
	}

	msgSize := t.order.Uint32(t.headLen)

	if msgSize > MsgMaxSize || msgSize < MsgHeadLen {
		return nil,fmt.Errorf("收到的数据长度非法:%d",size)
	}

	msgData := make([]byte,msgSize-MsgHeadLen)
	if _,err := io.ReadFull(reader,msgData); err != nil {
		return nil,err
	}

	msgType := t.order.Uint32(msgData)
	msgBody := msgData[MsgTypeLen:]

	return &LTVPacket{
		Len : msgSize,
		Type:  msgType,
		Value:  msgBody,
	},nil
}