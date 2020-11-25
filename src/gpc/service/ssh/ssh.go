package ssh

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"goPanel/src/constants"
	"goPanel/src/gpc/service"
	"goPanel/src/library/ssh"
	gossh "golang.org/x/crypto/ssh"
)

type TcpSsh struct {
	SshRead  chan []byte
	SshWrite chan []byte
}

func NewTcpSsh() *TcpSsh {
	return &TcpSsh{
		SshRead:  make(chan []byte, 1024),
		SshWrite: make(chan []byte, 1024),
	}
}

// 连接ssh
func (c *TcpSsh) SshConn(host, username, passwd string, port int, cols, rows uint32) (*ssh.Ssh, gossh.Channel, error) {
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

// 读ssh数据写入到tcp中
func (c *TcpSsh) SshReadBySocketWrite(tcpWrite chan []byte) {
	defer func() {
		log.Error(recover())
	}()

	for {
		select {
		case msg, ok := <-c.SshRead:
			if !ok {
				log.Error(ok)
			}
			wsMess := &service.Message{
				Event: constants.WS_EVENT_DATA,
				Data:  string(msg),
				Code:  constants.SUCCESS,
			}

			wsMessJson, _ := json.Marshal(wsMess)
			tcpWrite <- wsMessJson
		}
	}
}

// 读tcp数据写ssh数据
func (c *TcpSsh) ReadSocketBySshWrite(tcpRead chan []byte) {
	defer func() {
		log.Error(recover())
	}()

	for {
		select {
		case sw := <-tcpRead:
			c.SshWrite <- sw
		}
	}
}
