package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"goPanel/src/panel/constants"
	"goPanel/src/panel/library/ssh"
	"goPanel/src/panel/services"
	gossh "golang.org/x/crypto/ssh"
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
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type WsInitData struct {
	Host  string `json:"host"`
	Cols  uint32 `json:"cols"`
	Rows  uint32 `json:"rows"`
	Token string `json:"token"`
}

func NewWsController() *WsController {
	return &WsController{
		WsInit:      new(WsInitData),
		userService: new(services.UserService),
		initializer: make(chan bool, 1),
		WsRead:      make(chan []byte, 5120),
		WsWrite:     make(chan []byte, 10240),
		SshRead:     make(chan []byte, 10240),
		SshWrite:    make(chan []byte, 5120),
	}
}

func (c *WsController) Ssh(g *gin.Context) {
	var ch chan bool
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
	c.Ws = ws

	// 读ws客户端数据
	go c.wsRead()
	// 写ws客户端数据
	go c.wsWrite()

	select {
	case <-c.initializer:
		sh, sshChannel, err := c.sshConn(c.WsInit.Host, "fengxiao", "ZpB123", 22, c.WsInit.Cols, c.WsInit.Rows)
		if err != nil {
			log.Error(err)
			ret := &WsMessageData{
				Event: constants.WS_EVENT_ERR,
				Data:  constants.SSH_CONNECTION_FAILED_MSG,
			}
			retJson, _ := json.Marshal(ret)

			c.WsWrite <- retJson
		}
		defer func() {
			if err := sshChannel.Close(); err != nil {
				log.Error(err)
				return
			}
		}()

		// 转换为ws和ssh所识别的数据
		go c.sshRead()
		go c.sshWrite()

		go sh.Read(sshChannel, c.SshRead)
		go sh.Write(sshChannel, c.SshWrite)
	}

	defer func() {
		ch <- true
	}()

	<-ch
}

// 读ws数据
func (c *WsController) wsRead() {
	defer func() {
		recover()
		c.Ws.Close()
	}()

	for {
		mt, message, err := c.Ws.ReadMessage()
		// 其他错误，如果是 1001 和 1000 就不打印日志
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
			log.Info("ReadMessage other remote:%v error: %v \n", c.Ws.RemoteAddr(), err)
			return
		}

		if mt == websocket.BinaryMessage {
			reqMess := new(WsMessageData)
			_ = json.Unmarshal(message, &reqMess)

			switch reqMess.Event {
			case constants.WS_EVENT_INIT:
				reqDataJson, _ := json.Marshal(reqMess.Data)
				_ = json.Unmarshal(reqDataJson, &c.WsInit)

				// 验证登录情况
				if state, msg, _ := c.userService.IsUserLogin(c.WsInit.Token); !state {
					ret := &WsMessageData{
						Event: constants.WS_EVENT_ERR,
						Data:  msg,
					}
					retJson, _ := json.Marshal(ret)

					c.WsWrite <- retJson
				} else {
					c.initializer <- true
				}

				break
			case constants.WS_EVENT_DATA:
				c.WsRead <- []byte(reqMess.Data.(string))
				break
			}
		}
	}
}

// 写ws数据
func (c *WsController) wsWrite() {
	defer func() {
		recover()
		c.Ws.Close()
	}()

	for {
		select {
		case message := <-c.WsWrite:
			err := c.Ws.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Info(err)
				return
			}
		}
	}
}

// 连接ssh
func (c *WsController) sshConn(host, username, passwd string, port int, cols, rows uint32) (*ssh.Ssh, gossh.Channel, error) {
	sh := ssh.NewSsh(host, username, passwd, port)
	sshChannel, err := sh.RunShell(ssh.TermConfig{
		Cols: cols,
		Rows: rows,
	})
	if err != nil {
		log.Info(err)
		return nil, nil, err
	}

	return sh, sshChannel, nil
}

// 读ssh数据写入到ws中
func (c *WsController) sshRead() {
	for {
		select {
		case msg := <-c.SshRead:
			wsMess := &WsMessageData{
				Event: constants.WS_EVENT_DATA,
				Data:  string(msg),
			}

			wsMessJson, _ := json.Marshal(wsMess)
			c.WsWrite <- wsMessJson
		}
	}
}

// 读ws数据写ssh数据
func (c *WsController) sshWrite() {
	for {
		select {
		case msg := <-c.WsRead:
			c.SshWrite <- msg
		}
	}
}
