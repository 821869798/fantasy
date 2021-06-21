package api

type MsgHandle interface {
	TriggerEvent(interface{})
}