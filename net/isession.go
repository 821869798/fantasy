package net

type Session interface {
	Sid() uint64
	Send(msg interface{})
	Close()
	IsClose() bool
}