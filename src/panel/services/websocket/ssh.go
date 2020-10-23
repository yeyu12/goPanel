package websocket

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"goPanel/src/panel/constants"
	"goPanel/src/panel/library/ssh"
	gossh "golang.org/x/crypto/ssh"
)

type wsSsh struct {
	SshRead  chan []byte
	SshWrite chan []byte
}

// 连接ssh
func (c *wsSsh) sshConn(host, username, passwd string, port int, cols, rows uint32) (*ssh.Ssh, gossh.Channel, error) {
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
func (c *wsSsh) SshReadByWsWrite(wsWrite chan []byte) {
	defer func() {
		log.Error(recover())
	}()

	for {
		select {
		case msg, ok := <-c.SshRead:
			if !ok {
				log.Error(ok)
			}
			wsMess := &Message{
				Event: constants.WS_EVENT_DATA,
				Data:  string(msg),
				Code:  constants.SUCCESS,
			}

			wsMessJson, _ := json.Marshal(wsMess)
			wsWrite <- wsMessJson
		}
	}

	//CLOSE:
}

// 读ws数据写ssh数据
func (c *wsSsh) ReadWsBySshWrite(wsRead chan []byte) {
	defer func() {
		log.Error(recover())
	}()

	for {
		select {
		case msg := <-wsRead:
			c.SshWrite <- msg
		}
	}
}
