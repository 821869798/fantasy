package event

import "github.com/821869798/fantasy/net/api"

type SessionMsg struct {
	Session api.ISession
	Msg     interface{}
}

type SessionAdd struct {
	Session api.ISession
}

type SessionRemove struct {
	Session api.ISession
}
