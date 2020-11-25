package service

var ControlAddr string

type RequestWsMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	Code  int32       `json:"code"`
}
