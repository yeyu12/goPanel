package ssh

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
	"goPanel/src/constants"
	"goPanel/src/gpc/config"
	"goPanel/src/gpc/service"
	"goPanel/src/library/ssh"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"net"
)

type RelayClient struct {
	Conn    *net.TCPConn
	Read    chan []byte
	Write   chan []byte
	sh      *ssh.Ssh
	sshChan gossh.Channel
}

func NewRelayClient() *RelayClient {
	return &RelayClient{
		Read:  make(chan []byte, 1024),
		Write: make(chan []byte, 1024),
		Conn:  new(net.TCPConn),
	}
}

func (r *RelayClient) RelayConn(addr string, flag uint, cols, rows uint32) error {
	relayConn, err := common.ConnTcp(addr)
	if err != nil {
		return err
	}

	r.Conn = relayConn

	go r.relayClientRead()
	go r.relayClientWrite()

	switch flag {
	case constants.CLIENT_SHELL_TYPE:
		tcpSsh, err := r.connSsh(cols, rows)
		if err != nil {
			return err
		}

		// tcp和ssh交换数据
		go r.relayClientReadTcpWriteSsh(tcpSsh.SshWrite)
		go r.relayClientReadSshWriteTcp(tcpSsh.SshRead)

		go r.sh.Read(r.sshChan, tcpSsh.SshRead)
		go r.sh.Write(r.sshChan, tcpSsh.SshWrite)

		break
	}

	return nil
}

func (r *RelayClient) relayClientReadTcpWriteSsh(sshWrite chan []byte) {
	for {
		select {
		case sw := <-r.Read:
			sshWrite <- sw
		}
	}
}

func (r *RelayClient) relayClientReadSshWriteTcp(sshRead chan []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	for {
		select {
		case w := <-sshRead:
			r.Write <- w
		}
	}
}

func (r *RelayClient) connSsh(cols, rows uint32) (*TcpSsh, error) {
	var err error
	host := "127.0.0.1"
	tcpSsh := NewTcpSsh()
	r.sh, r.sshChan, err = tcpSsh.SshConn(host, config.Conf.Ssh.Username, config.Conf.Ssh.Password, config.Conf.Ssh.Port, cols, rows)
	if err != nil {
		m := service.Message{
			Event: constants.WS_EVENT_ERR,
			Data:  constants.SSH_CONNECTION_FAILED_MSG,
			Code:  constants.SSH_CONNECTION_FAILED,
		}

		resJson, _ := json.Marshal(m)
		_, err = r.Conn.Write(resJson)
		if err != nil {
			log.Error(err)
			r.closeRelay()
		}

		return nil, err
	}

	return tcpSsh, nil
}

func (r *RelayClient) relayClientRead() {
	for {
		data := make([]byte, 1024)
		size, err := r.Conn.Read(data)
		if (err != nil) || (err == io.EOF) {
			if err != io.EOF {
				log.Error(err)
			}
			r.closeRelay()

			break
		}

		r.Read <- data[:size]
		log.Info(string(data[:size]))
	}
}

func (r *RelayClient) relayClientWrite() {
	for {
		select {
		case wr := <-r.Write:
			_, err := r.Conn.Write(wr)
			if err != nil {
				log.Error(err)
				r.closeRelay()
			}
		}
	}
}

func (r *RelayClient) closeRelay() {
	if r.Conn != nil {
		r.Conn.Close()
	}
	if r.sshChan != nil {
		_ = r.sshChan.Close()
	}
}
