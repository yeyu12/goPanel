package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"goPanel/src/panel/common"
	"goPanel/src/panel/library/ssh"
	"net/http"
	"time"
)

type WsController struct {
	BaseController
}

func NewWsController() *WsController {
	return &WsController{}
}

func (c *WsController) Ssh(g *gin.Context) {
	ws, err := (&websocket.Upgrader{
		HandshakeTimeout: time.Duration(time.Second * 30),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		return
	}
	defer func() {
		_ = ws.Close()
	}()

	cols, _ := common.StringUtils(g.Param("cols")).Uint32()
	rows, _ := common.StringUtils(g.Param("rows")).Uint32()
	host := g.Param("host")

	// 通过ip获取相关ssh客户端数据
	sh := ssh.NewSsh(host, "fengxiao", "ZpB123", 22)
	sshChannel, err := sh.RunShell(ssh.TermConfig{
		Cols: cols,
		Rows: rows,
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := sshChannel.Close(); err != nil {
			log.Error(err)
		}
	}()

	wsRead := make(chan []byte, 5120)
	wsWrite := make(chan []byte, 10240)
	var ch chan bool

	go sh.Read(sshChannel, wsWrite)
	go sh.Write(sshChannel, wsRead)

	// 读ws客户端数据
	go func() {
		defer func() {
			ch <- true
		}()

		for {
			mt, message, err := ws.ReadMessage()
			// 其他错误，如果是 1001 和 1000 就不打印日志
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Debug("ReadMessage other remote:%v error: %v \n", ws.RemoteAddr(), err)
				return
			}

			if mt == websocket.TextMessage {
				wsRead <- message
			}
		}
	}()

	// 写ws客户端数据
	go func() {
		defer func() {
			ch <- true
		}()

		for {
			select {
			case message := <-wsWrite:
				err = ws.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Debug(err)
					return
				}
			}
		}
	}()

	<-ch
}
