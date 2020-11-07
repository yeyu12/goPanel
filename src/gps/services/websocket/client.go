package websocket

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"goPanel/src/gps/common"
	"goPanel/src/gps/constants"
	core "goPanel/src/gps/core/database"
	"goPanel/src/gps/services"
	"time"
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
		time.Sleep(time.Microsecond * 100)
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
					c.wsWriteErr(code, msg)
					log.Info(msg)

					return
				}

				switch reqMess.Type {
				case CLIENT_SHELL_TYPE:
					sshInitData := new(ShellInit)
					_ = json.Unmarshal(baseInitJson, &sshInitData)

					// 查询相关数据
					service := new(services.MachineService)
					machineData := service.IdByDetails(core.Db, sshInitData.Id)
					if machineData.Id == 0 {
						c.wsWriteErr(constants.MACHINE_ID, constants.MACHINE_ID_MSG)
						log.Info(constants.MACHINE_ID_MSG)

						return
					}
					if sshInitData.Passwd == "" {
						c.wsWriteErr(constants.MACHINE_PASSWD_NOT_NILL, constants.MACHINE_PASSWD_NOT_NILL_MSG)
						log.Info(constants.MACHINE_PASSWD_NOT_NILL_MSG)

						return
					}

					sDec, _ := base64.StdEncoding.DecodeString(sshInitData.Passwd)
					passwd, err := common.RsaDecrypt(sDec, common.GetRsaFilePath()+"private.pem")
					if err != nil {
						c.wsWriteErr(constants.MACHINE_PASSWD_DECODE_FAIL, constants.MACHINE_PASSWD_DECODE_FAIL_MSG)
						log.Info(err)

						return
					}

					sh, sshChannel, err := c.wsShell.sshConn(machineData.Host, machineData.User, string(passwd), machineData.Port, sshInitData.Cols, sshInitData.Rows)
					if err != nil {
						c.wsWriteErr(constants.SSH_CONNECTION_FAILED, constants.SSH_CONNECTION_FAILED_MSG)
						log.Info(err)

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
					go c.wsShell.sshReadByWsWrite(c.Send)
					go c.wsShell.readWsBySshWrite(c.wsRead)
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
