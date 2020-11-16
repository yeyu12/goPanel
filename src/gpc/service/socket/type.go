package socket

var isReconnControlTcp = false

type RequestWsMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	Code  int32       `json:"code"`
}
