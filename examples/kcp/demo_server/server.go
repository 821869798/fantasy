package main

import (
	"github.com/821869798/fantasy/net/kcp"
	"github.com/821869798/fantasy/net/network"
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
	case *network.SessionMsg:
		m := e.(*network.SessionMsg)
		p, ok := m.Msg.(*packet.LTVPacket)
		if ok {
			slog.Infof("MsgHandle recv client msg:%s", string(p.Value))
		}
		_ = m.Session.Send(p)
	case *network.SessionAdd:
		m := e.(*network.SessionAdd)
		slog.Infof("Session connected :%v", m.Session.RemoteAddr())
	case *network.SessionRemove:
		m := e.(*network.SessionRemove)
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
