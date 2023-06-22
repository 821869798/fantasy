package api

// MsgCodec 消息收发编解码
type MsgCodec interface {
	OnSendMsg(s Session, msg interface{}) error
	OnRecvMsg(s Session) (interface{}, error)
}
