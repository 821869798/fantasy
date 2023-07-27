package ws

import (
	"encoding/binary"
	"errors"
	"github.com/821869798/fantasy/net/network"
	"github.com/gorilla/websocket"
)

const (
	MsgTypeLen = 4
	MsgHeadLen = MsgTypeLen
	MsgMaxSize = 65536
)

type WsPacket struct {
	Type  uint32
	Value []byte
}

func NewWsPacket(t uint32, v []byte) *WsPacket {
	return &WsPacket{
		Type:  t,
		Value: v,
	}
}

type wsMsgCodec struct {
	order binary.ByteOrder
}

func NewWsMsgCodec(order binary.ByteOrder) network.IMsgCodec {
	t := &wsMsgCodec{
		order: order,
	}
	return t
}

func (t *wsMsgCodec) OnSendMsg(s network.ISession, msg interface{}) error {
	writer, ok := s.Raw().(*websocket.Conn)
	if !ok || writer == nil {
		return errors.New("error to send msg,not a websocket conn")
	}

	packet, ok := msg.(*WsPacket)
	if !ok {
		return nil
	}

	rawData := make([]byte, len(packet.Value)+MsgHeadLen)

	t.order.PutUint32(rawData, packet.Type)
	copy(rawData[MsgHeadLen:], packet.Value)

	err := writer.WriteMessage(websocket.BinaryMessage, rawData)

	return err
}
func (t *wsMsgCodec) OnRecvMsg(s network.ISession) (interface{}, error) {
	reader, ok := s.Raw().(*websocket.Conn)
	if !ok {
		return nil, errors.New("error to read,not a websocket conn")
	}

	wsMsgType, msgData, err := reader.ReadMessage()
	_ = wsMsgType
	if err != nil {
		return nil, err
	}

	msgType := t.order.Uint32(msgData)
	msgBody := msgData[MsgHeadLen:]

	packet := NewWsPacket(msgType, msgBody)

	return packet, nil
}
