package ws

// WsStartOpt 启动参数
type WsStartOpt struct {
	MaxConns        int
	ReadBufferSize  int
	WriteBufferSize int
	HandlePath      string
}

func NewWsStartOpt() *WsStartOpt {
	o := &WsStartOpt{
		MaxConns:        20000,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		HandlePath:      "/",
	}
	return o
}

func (o *WsStartOpt) SetMaxConns(maxConns int) *WsStartOpt {
	o.MaxConns = maxConns
	return o
}
