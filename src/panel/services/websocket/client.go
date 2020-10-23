package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"goPanel/src/panel/constants"
	"goPanel/src/panel/services"
)

type Client struct {
	UID        string
	Socket     *websocket.Conn
	Send       chan []byte
	wsRead     chan []byte
	ClientType int
	wsShell    *wsSsh
}

var userService = new(services.UserService)

func NewWsShell(uid string, socket *websocket.Conn) *Client {
	return &Client{
		UID:        uid,
		Socket:     socket,
		Send:       make(chan []byte, 1024),
		ClientType: 0,
		wsShell: &wsSsh{
			SshRead:  make(chan []byte, 1024),
			SshWrite: make(chan []byte, 1024),
		},
		wsRead: make(chan []byte, 1024),
	}
}

func (c *Client) Read() {
	defer func() {
		recover()
		WsManager.UnRegister <- c
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

			switch reqMess.Event {
			case constants.WS_EVENT_INIT:
				// 判断登录情况
				baseInitData := new(BaseInit)
				baseInitJson, _ := json.Marshal(reqMess.Data)
				_ = json.Unmarshal(baseInitJson, &baseInitData)

				if state, msg, code := userService.IsUserLogin(baseInitData.Token); !state {
					ret := &Message{
						Event: constants.WS_EVENT_ERR,
						Data:  msg,
						Code:  code,
					}
					retJson, _ := json.Marshal(ret)

					c.Send <- retJson
					return
				}

				switch reqMess.Type {
				case CLIENT_SHELL_TYPE:
					sshInitData := new(ShellInit)
					_ = json.Unmarshal(baseInitJson, &sshInitData)

					sh, sshChannel, err := c.wsShell.sshConn(sshInitData.Host, "yeyu", "ZpB123", 22, sshInitData.Cols, sshInitData.Rows)
					if err != nil {
						ret := &Message{
							Event: constants.WS_EVENT_ERR,
							Data:  constants.SSH_CONNECTION_FAILED_MSG,
							Code:  constants.SSH_CONNECTION_FAILED,
						}

						retJson, _ := json.Marshal(ret)
						c.Send <- retJson

						log.Error(err)
						return
					}
					defer func() {
						if err := sshChannel.Close(); err != nil {
							log.Error(err)
							return
						}
					}()

					// 读ws  转换sshWrite    sshWrite写入通道
					// 写ws  转换sshRead     sshRead读通道
					// 转换为ws和ssh所识别的数据
					go c.wsShell.SshReadByWsWrite(c.Send)
					go c.wsShell.ReadWsBySshWrite(c.wsRead)
					go sh.Read(sshChannel, c.wsShell.SshRead)
					go sh.Write(sshChannel, c.wsShell.SshWrite)

					break
				}

				break
			case constants.WS_EVENT_DATA:
				switch reqMess.Type {
				case CLIENT_SHELL_TYPE:
					c.wsRead <- []byte(reqMess.Data.(string))
					break
				}

				break
			}
		}
	}
}

func (c *Client) Write() {
	defer func() {
		recover()
		WsManager.UnRegister <- c
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
