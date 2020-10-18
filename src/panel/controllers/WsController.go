package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}
