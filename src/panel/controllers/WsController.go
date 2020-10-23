package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"goPanel/src/panel/constants"
	"goPanel/src/panel/library/ssh"
	websocket2 "goPanel/src/panel/library/websocket"
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
	Type  int         `json:"type"`
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

func (c *WsController) SshN(g *gin.Context) {
	ws, err := (&websocket.Upgrader{
		HandshakeTimeout: time.Duration(time.Second * 30),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		log.Error(err)
		return
	}

	client := &websocket2.Client{
		ID:         uuid.NewV4().String(),
		Socket:     ws,
		Send:       make(chan []byte),
		ClientType: websocket2.CLIENT_SHELL_TYPE,
	}

	websocket2.Manager.Register <- client

	go client.Read()
	go client.Write()
}

func (c *WsController) Ssh(g *gin.Context) {
	var ch chan bool
	var closeSsh = make(chan bool, 4)

	ws, err := (&websocket.Upgrader{
		HandshakeTimeout: time.Duration(time.Second * 30),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer func() {
		_ = ws.Close()
	}()
	c.Ws = ws

	// 读ws客户端数据
	go c.wsRead(closeSsh)
	// 写ws客户端数据
	go c.wsWrite(closeSsh)

	select {
	case <-c.initializer:
		sh, sshChannel, err := c.sshConn(c.WsInit.Host, "fengxiao", "ZpB123", 22, c.WsInit.Cols, c.WsInit.Rows)
		if err != nil {
			ret := &WsMessageData{
				Event: constants.WS_EVENT_ERR,
				Data:  constants.SSH_CONNECTION_FAILED_MSG,
			}
			retJson, _ := json.Marshal(ret)

			c.WsWrite <- retJson
		}
		defer func() {
			if err := sshChannel.Close(); err != nil {
				return
			}
		}()

		// 转换为ws和ssh所识别的数据
		go c.sshRead(closeSsh)
		go c.sshWrite(closeSsh)

		go sh.Read(closeSsh, sshChannel, c.SshRead)
		go sh.Write(closeSsh, sshChannel, c.SshWrite)
	}

	defer func() {
		ch <- true
	}()

	<-ch
}

// 读ws数据
func (c *WsController) wsRead(closeChan chan bool) {
	defer func() {
		recover()
		close(closeChan)
		_ = c.Ws.Close()
	}()

	for {
		mt, message, err := c.Ws.ReadMessage()
		// 其他错误，如果是 1001 和 1000 就不打印日志
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
			log.Info("ReadMessage other remote:%v error: %v \n", c.Ws.RemoteAddr(), err)
			return
		}

		log.Error(&c.Ws)

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
func (c *WsController) wsWrite(closeChan chan bool) {
	defer func() {
		recover()
		close(closeChan)
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
		return nil, nil, err
	}

	return sh, sshChannel, nil
}

// 读ssh数据写入到ws中
func (c *WsController) sshRead(closeChan chan bool) {
	for {
		select {
		case msg := <-c.SshRead:
			wsMess := &WsMessageData{
				Event: constants.WS_EVENT_DATA,
				Data:  string(msg),
			}

			wsMessJson, _ := json.Marshal(wsMess)
			c.WsWrite <- wsMessJson
		case <-closeChan:
			goto CLOSE
		}
	}

CLOSE:
}

// 读ws数据写ssh数据
func (c *WsController) sshWrite(closeChan chan bool) {
	for {
		select {
		case msg := <-c.WsRead:
			c.SshWrite <- msg
		case <-closeChan:
			goto CLOSE
		}
	}

CLOSE:
}
