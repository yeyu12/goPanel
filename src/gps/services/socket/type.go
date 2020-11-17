package socket

import (
	"goPanel/src/gps/config"
	"strconv"
)

var (
	controlAddr = ":" + strconv.Itoa(config.Conf.App.ControlPort) // 控制端口
)

var ServerWsManager = ServerWebsocketManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	UnRegister: make(chan *Client),
	Clients:    make(map[*Client]bool),
}

var ControlManager = &ControlTcpManager{
	Client:         make(map[*Control]bool),
	Broadcast:      make(chan []byte, 1024),
	Register:       make(chan *Control),
	UnRegister:     make(chan *Control),
	relayStartPort: config.Conf.App.RelayStartPort,
}

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
