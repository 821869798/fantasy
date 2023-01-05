package tcp

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
	Len   uint32
	Type  uint32
	Value []byte
}

func NewTcpPacket(t uint32, v []byte) *LTVPacket {
	return &LTVPacket{
		Len:   uint32(len(v)) + MsgHeadLen,
		Type:  t,
		Value: v,
	}
}

type tcpTransmitter struct {
	order binary.ByteOrder
}

func NewTcpTransmitter(order binary.ByteOrder) api.MsgTransmitter {
	t := &tcpTransmitter{
		order: order,
	}
	return t
}

func (t *tcpTransmitter) OnSendMsg(s api.Session, msg interface{}) error {
	writer, ok := s.Raw().(io.Writer)
	if !ok || writer == nil {
		return nil
	}

	packet, ok := msg.(*LTVPacket)
	if !ok {
		return nil
	}

	rawData := make([]byte, packet.Len)

	t.order.PutUint32(rawData, packet.Len)
	t.order.PutUint32(rawData[MsgSizeLen:], packet.Type)
	copy(rawData[MsgHeadLen:], packet.Value)

	_, err := writer.Write(rawData)

	return err
}
func (t *tcpTransmitter) OnRecvMsg(s api.Session) (interface{}, error) {
	reader, ok := s.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	headLen := make([]byte, MsgSizeLen)
	if _, err := io.ReadFull(reader, headLen); err != nil {
		return nil, err
	}

	msgSize := t.order.Uint32(headLen)

	if msgSize > MsgMaxSize || msgSize < MsgHeadLen {
		return nil, fmt.Errorf("recv packet length error:%d", msgSize)
	}

	msgData := make([]byte, msgSize-MsgSizeLen)
	if _, err := io.ReadFull(reader, msgData); err != nil {
		return nil, err
	}

	msgType := t.order.Uint32(msgData)
	msgBody := msgData[MsgTypeLen:]

	packet := NewTcpPacket(msgType, msgBody)

	return packet, nil
}
