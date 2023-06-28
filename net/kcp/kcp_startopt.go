package kcp

import "time"

// KcpStartOpt 启动参数
type KcpStartOpt struct {
	MaxConns     int
	SendChanSize uint32
	Timeout      time.Duration
}

func NewKcpStartOpt() *KcpStartOpt {
	o := &KcpStartOpt{
		MaxConns:     20000,
		SendChanSize: 1024,
		Timeout:      15 * time.Second,
	}
	return o
}

func (o *KcpStartOpt) SetMaxConns(maxConns int) *KcpStartOpt {
	o.MaxConns = maxConns
	return o
}

func (o *KcpStartOpt) SetSendChanSize(chanSize uint32) *KcpStartOpt {
	o.SendChanSize = chanSize
	return o
}
