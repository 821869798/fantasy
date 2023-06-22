package main

import (
	"github.com/821869798/fantasy/net/event"
	"github.com/821869798/fantasy/net/tcp"
	"github.com/gookit/slog"
	"os"
	"os/signal"
	"syscall"
)

type MsgHandle struct {
}

func (m *MsgHandle) TriggerEvent(e interface{}) {
	switch e.(type) {
	case *event.SessionMsg:
		m := e.(*event.SessionMsg)
		packet, ok := m.Msg.(*tcp.LTVPacket)
		if ok {
			slog.Infof("MsgHandle recv client msg:%s", string(packet.Value))
		}
		m.Session.Send(packet)
	case *event.SessionAdd:
		m := e.(*event.SessionAdd)
		slog.Infof("Session connected :%v", m.Session.RemoteAddr())
	case *event.SessionRemove:
		m := e.(*event.SessionRemove)
		slog.Infof("Session disconnected :%v", m.Session.RemoteAddr())
	}
}

func main() {

	slog.SetLogLevel(slog.DebugLevel)

	slog.Infof("server start...")
	a := tcp.NewTcpAcceptor(":7801", &MsgHandle{}, nil, nil)
	a.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
}
