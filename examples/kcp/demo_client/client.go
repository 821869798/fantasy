package main

import (
	"fmt"
	"github.com/821869798/fantasy/net/kcp"
	"github.com/821869798/fantasy/net/network"
	"github.com/821869798/fantasy/net/packet"
	"github.com/gookit/slog"
)

type MsgHandle struct {
}

func (m *MsgHandle) TriggerEvent(e interface{}) {
	switch m := e.(type) {
	case *network.SessionMsg:
		p, ok := m.Msg.(*packet.LTVPacket)
		if ok {
			slog.Infof("MsgHandle recv server msg:%s", string(p.Value))
		}
	}
}

func main() {

	slog.SetLogLevel(slog.DebugLevel)

	c := kcp.NewKcpConnector("127.0.0.1:7801", &MsgHandle{}, nil, nil)
	c.Start()

	var input string
	for true {
		_, err := fmt.Scanln(&input)
		if err != nil {
			slog.Errorf("%v", err)
			return
		}

		p := packet.NewLTVPacket(0, []byte(input))

		err = c.Session().Send(p)
		if err != nil {
			slog.Errorf("client send error:%v", err)
		}
	}
}
