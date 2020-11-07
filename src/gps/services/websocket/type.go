package websocket

const (
	CLIENT_SHELL_TYPE = iota
)

type Message struct {
	Type  int         `json:"type"`
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	Code  int32       `json:"code"`
}

type BaseInit struct {
	Token string `json:"token"`
}

type ShellInit struct {
	BaseInit
	Id     int64  `json:"id"`
	Passwd string `json:"passwd"`
	Cols   uint32 `json:"cols"`
	Rows   uint32 `json:"rows"`
}
