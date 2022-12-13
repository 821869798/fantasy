package api

import "net"

type Session interface {
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
