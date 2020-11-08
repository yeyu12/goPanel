package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"goPanel/src/gps/services"
	panel_ws "goPanel/src/gps/services/websocket"
	"goPanel/src/gps/services/ws"
	"net/http"
	"time"
)

type WsController struct {
	BaseController
	Ws          *websocket.Conn
	WsInit      *WsInitData
	userService *services.UserService
	initializer chan bool
	WsRead      chan []byte
	WsWrite     chan []byte
	SshRead     chan []byte
	SshWrite    chan []byte
}

type WsMessageData struct {
	Type  int         `json:"type"`
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
}

type WsInitData struct {
	Host  string `json:"host"`
	Cols  uint32 `json:"cols"`
	Rows  uint32 `json:"rows"`
	Token string `json:"token"`
}

func NewWsController() *WsController {
	return &WsController{
		Ws:          new(websocket.Conn),
		WsInit:      new(WsInitData),
		userService: new(services.UserService),
		initializer: make(chan bool, 1),
		WsRead:      make(chan []byte, 1024),
		WsWrite:     make(chan []byte, 1024),
		SshRead:     make(chan []byte, 1024),
		SshWrite:    make(chan []byte, 1024),
	}
}

func (c *WsController) Ssh(g *gin.Context) {
	wsConn, err := (&websocket.Upgrader{
		HandshakeTimeout: time.Duration(time.Second * 30),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		log.Error(err)
		return
	}

	client := panel_ws.NewWsShell(uuid.NewV4().String(), wsConn)
	panel_ws.WsManager.Register <- client

	go client.Read()
	go client.Write()
}

func (c *WsController) Conn(g *gin.Context) {
	wsConn, err := (&websocket.Upgrader{
		HandshakeTimeout: time.Duration(time.Second * 30),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		log.Error(err)
		return
	}

	client := ws.NewWs(uuid.NewV4().String(), wsConn)
	ws.WsManager.Register <- client

	go client.Read()
	go client.Write()
}
