package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"goPanel/src/constants"
	"goPanel/src/gps/services"
	"net"
)

type Client struct {
	UID           string
	Socket        *websocket.Conn
	Send          chan []byte
	wsRead        chan []byte
	ClientType    int
	RelayListener *net.TCPListener
}

var userService = new(services.UserService)

func NewClientWs(uid string, socket *websocket.Conn) *Client {
	return &Client{
		UID:        uid,
		Socket:     socket,
		Send:       make(chan []byte, 1024),
		ClientType: 0,
		wsRead:     make(chan []byte, 1024),
	}
}

func (c *Client) Read() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}

		//time.Sleep(time.Microsecond * 100)
		ServerWsManager.UnRegister <- c
	}()

	for {
		mt, message, err := c.Socket.ReadMessage()
		// 其他错误，如果是 1001 和 1000 就不打印日志
		if websocket.IsUnexpectedCloseError(err,
			websocket.CloseGoingAway,
			websocket.CloseNormalClosure,
			websocket.CloseNoStatusReceived,
		) {
			log.Infof("ReadMessage other remote:%v error: %v \n", c.Socket.RemoteAddr(), err)
			return
		}

		if mt == websocket.BinaryMessage {
			reqMess := new(Message)
			_ = json.Unmarshal(message, &reqMess)

			c.handleWsMess(reqMess)
		}
	}
}

func (c *Client) Write() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}

		ServerWsManager.UnRegister <- c
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			msgJson, _ := json.Marshal(Message{
				Type:  0,
				Event: constants.WS_EVENT_DATA,
				Data:  string(message),
				Code:  constants.SUCCESS,
			})

			if err := c.Socket.WriteMessage(websocket.BinaryMessage, msgJson); err != nil {
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

func (c *Client) handleWsMess(req *Message) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	switch req.Event {
	case constants.WS_EVENT_INIT:
		// 验证相关数据，用户token，终端是否存在
		baseInitData := new(BaseInit)
		baseInitJson, _ := json.Marshal(req.Data)
		_ = json.Unmarshal(baseInitJson, &baseInitData)

		if state, msg, code := userService.IsUserLogin(baseInitData.Token); !state {
			c.wsWriteErr(code, msg)
			log.Info(msg)

			return
		}

		switch req.Type {
		case constants.CLIENT_SHELL_TYPE:
			sshInitData := new(ShellInit)
			_ = json.Unmarshal(baseInitJson, &sshInitData)

			// 查询是否有该连接存在
			// 存在：创建中继端，通知客户端创建一个连接，连接中继端，客户端返回终端密码，ws客户端直连中继端
			// 不存在：跳出
			control := ControlManager.FindClientIdByClientConn(sshInitData.Id)
			if control == nil {
				c.wsWriteErr(constants.CLIENT_NOT_FOND_FAIL, constants.CLIENT_NOT_FOND_MSG)
				return
			}

			// 创建中继端
			port := RelayPort()
			relayListener, err := CreateRelayConn(port, c.Send)
			if err != nil {
				log.Error(err)
				c.wsWriteErr(constants.CREATE_NOT_RELAY_FAIL, constants.CREATE_NOT_RELAY_MSG)
				return
			}
			c.RelayListener = relayListener

			// 通知客户端创建本地ssh，连接中间端
			reqMess := RequestWsMessage{
				Event: "sshConnectRelay",
				Data: map[string]interface{}{
					"port": port,
					"cols": sshInitData.Cols,
					"rows": sshInitData.Rows,
				},
			}
			reqMessJson, _ := json.Marshal(reqMess)
			control.Write <- reqMessJson

			break
		}

		break
	case constants.WS_EVENT_DATA:
		switch req.Type {
		case constants.CLIENT_SHELL_TYPE:
			log.Error(req.Data.(string))
			c.wsRead <- []byte(req.Data.(string))
			break
		}

		break
	}
}
