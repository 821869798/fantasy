package api

import "net"

type Session interface {
	// 原始的Socket连接
	Raw() interface{}
	// 链接的地址
	RemoteAddr() net.Addr
	//session Id
	Sid() uint64
	//发送消息
	Send(msg interface{}) error
	//关闭
	Close()
	//是否关闭
	IsClose() bool
}
