package main

import (
	"github.com/821869798/fantasy/net/event"
	"github.com/821869798/fantasy/net/kcp"
	"github.com/821869798/fantasy/net/packet"
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
		p, ok := m.Msg.(*packet.LTVPacket)
		if ok {
			slog.Infof("MsgHandle recv client msg:%s", string(p.Value))
		}
		_ = m.Session.Send(p)
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

	slog.Infof("kcp server start...")
	a := kcp.NewKcpAcceptor("127.0.0.1:7801", &MsgHandle{}, nil, nil)
	//public internet //a := tcp.NewKcpAcceptor(":7801", &MsgHandle{}, nil, nil)
	a.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
}
