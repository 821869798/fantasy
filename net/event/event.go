package event

import "github.com/821869798/fantasy/net/api"

type SessionMsg struct {
	Session api.Session
	Msg     interface{}
}

type SessionAdd struct {
	Session api.Session
}
