package api

type Session interface {
	// 原始的Socket连接
	Raw() interface{}
	//session Id
	Sid() uint64
	//发送消息
	Send(msg interface{}) error
	//关闭
	Close()
	//是否关闭
	IsClose() bool
}