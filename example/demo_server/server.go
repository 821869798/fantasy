package main

import (
	"github.com/821869798/fantasy/net/event"
	"github.com/821869798/fantasy/net/tcp"
	log "github.com/FishGoddess/logit"
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
			log.Info("MsgHandle recv client msg:%s", string(packet.Value))
		}
		m.Session.Send(packet)
	}
}

func main() {

	log.Me().SetLevel(log.DebugLevel)
	log.Me().NeedCaller(true)

	log.Info("server start...")
	a := tcp.NewTcpAcceptor(":7801", &MsgHandle{}, nil, nil)
	a.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
}
