package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"goPanel/src/gps/constants"
	"goPanel/src/gps/services"
	"time"
)

type Client struct {
	UID        string
	Socket     *websocket.Conn
	Send       chan []byte
	WsRead     chan []byte
	ClientType int
}

var userService = new(services.UserService)

func NewClientWs(uid string, socket *websocket.Conn) *Client {
	return &Client{
		UID:        uid,
		Socket:     socket,
		Send:       make(chan []byte, 1024),
		ClientType: 0,
		WsRead:     make(chan []byte, 1024),
	}
}

func (c *Client) Read() {
	defer func() {
		recover()
		time.Sleep(time.Microsecond * 100)
		ServerWsManager.UnRegister <- c
	}()

	for {
		mt, message, err := c.Socket.ReadMessage()
		// 其他错误，如果是 1001 和 1000 就不打印日志
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
			log.Infof("ReadMessage other remote:%v error: %v \n", c.Socket.RemoteAddr(), err)
			return
		}

		if mt == websocket.BinaryMessage {
			reqMess := new(Message)
			_ = json.Unmarshal(message, &reqMess)
		}
	}
}

func (c *Client) Write() {
	defer func() {
		recover()
		ServerWsManager.UnRegister <- c
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Socket.WriteMessage(websocket.BinaryMessage, message); err != nil {
				log.Error(err)
				return
			}
		}
	}
}

func (c *Client) wsWriteErr(code int32, msg string) {
	ret := &Message{
		Event: constants.WS_EVENT_ERR,
		Data:  msg,
		Code:  code,
	}
	retJson, _ := json.Marshal(ret)
	c.Send <- retJson

	return
}
