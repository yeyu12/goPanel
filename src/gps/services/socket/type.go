package socket

const (
	CLIENT_SHELL_TYPE = iota
)

var (
	serverAddr  = "0.0.0.0:10000" // 服务端口
	controlAddr = "0.0.0.0:10010" // 控制端口
	relayAddr   = "0.0.0.0:10086" // 中继端口
)

type Message struct {
	Type  int         `json:"type"`
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	Code  int32       `json:"code"`
}

type RequestWsMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
