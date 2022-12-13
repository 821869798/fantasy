package tcp

//启动参数
type TcpStartOpt struct {
	MaxConns     int
	SendChanSize uint32
}

func NewTcpStartOpt() *TcpStartOpt {
	o := &TcpStartOpt{
		MaxConns:     20000,
		SendChanSize: 1024,
	}
	return o
}

func (o *TcpStartOpt) SetMaxConns(maxConns int) *TcpStartOpt {
	o.MaxConns = maxConns
	return o
}

func (o *TcpStartOpt) SetSendChanSize(chanSize uint32) *TcpStartOpt {
	o.SendChanSize = chanSize
	return o
}
