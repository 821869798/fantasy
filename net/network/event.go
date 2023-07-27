package network

type SessionMsg struct {
	Session ISession
	Msg     interface{}
}

type SessionAdd struct {
	Session ISession
}

type SessionRemove struct {
	Session ISession
}
