package event

import "github.com/821869798/fantasy/net/api"

type EventSessionMsg struct {
	Session api.Session
	Msg interface{}
}

type EventSessionAdd struct {
	Session api.Session
}


