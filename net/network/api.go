package network

import "net"

type ISession interface {
	// Raw 原始的Socket连接
	Raw() interface{}
	// RemoteAddr 链接的地址
	RemoteAddr() net.Addr
	// Sid session Id
	Sid() uint64
	// Send 发送消息
	Send(msg interface{}) error
	// Close 关闭
	Close()
	// IsClose 是否关闭
	IsClose() bool
}

type INetwork interface {
	CreateListener(addr string) (net.Listener, bool)
	Dial(addr string) (net.Conn, bool)
	SessionAdapter() ISessionAdapter
}

type ISessionAdapter interface {
	Name() string
	SendChanSize() uint32
	RemoteAddr(rawConn interface{}) net.Addr
	CloseConn(rawConn interface{}) error
	Handle() IMsgHandle
	SendMsg(s ISession, msg interface{}) error
	RecvMsg(s ISession) (interface{}, error)
}

// IMsgHandle 消息接受处理器
type IMsgHandle interface {
	TriggerEvent(interface{})
}

// IMsgCodec 消息收发编解码
type IMsgCodec interface {
	OnSendMsg(s ISession, msg interface{}) error
	OnRecvMsg(s ISession) (interface{}, error)
}
