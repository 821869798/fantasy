package event

import "github.com/821869798/fantasy/net/api"

type Event interface {
	Session() api.Session
	Data() interface{}
}

type MsgEventSessionAdd struct {
	s api.Session
}

func NewMsgEventSessionAdd(s api.Session) MsgEventSessionAdd{
	e := &MsgEventSessionAdd{
		s : s,
	}
	return e
}

