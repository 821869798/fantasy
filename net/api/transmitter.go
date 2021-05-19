package api

// 消息收发器
type MsgTransmitter interface {
	OnSendMsg(s Session,msg interface{}) error
	OnRecvMsg(s Session) (interface{},error)
}