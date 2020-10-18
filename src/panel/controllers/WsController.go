package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"goPanel/src/panel/common"
	"goPanel/src/panel/library/ssh"
	"log"
	"net/http"
)

func Ssh(c *gin.Context) {
	ws, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	cols, _ := common.StringUtils(c.Param("cols")).Uint32()
	rows, _ := common.StringUtils(c.Param("rows")).Uint32()
	host := c.Param("host")

	// 通过ip获取相关ssh客户端数据
	sh := ssh.NewSsh(host, "yeyu", "ZpB123", 22)
	sshChannel, err := sh.RunShell(ssh.TermConfig{
		Cols: cols,
		Rows: rows,
	})

	wsRead := make(chan []byte, 1024)
	wsWrite := make(chan []byte, 1024)
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
			if err != nil {
				log.Fatal(err)
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
					log.Fatal(err)
				}
			}
		}
	}()

	<-ch
}
