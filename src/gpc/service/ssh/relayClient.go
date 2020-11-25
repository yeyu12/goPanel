package ssh

import (
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
	"io"
	"net"
)

type RelayClient struct {
	Conn  *net.TCPConn
	Read  chan []byte
	Write chan []byte
}

func NewRelayClient() *RelayClient {
	return &RelayClient{
		Read:  make(chan []byte, 1024),
		Write: make(chan []byte, 1024),
		Conn:  new(net.TCPConn),
	}
}

func (r *RelayClient) RelayConn(addr string) error {
	// 连接中继端
	relayConn, err := common.ConnTcp(addr)
	if err != nil {
		return err
	}

	r.Conn = relayConn

	go r.relayClientRead()
	go r.relayClientWrite()

	return nil
}

func (r *RelayClient) RelayClientReadTcpWriteSsh(sshWrite chan []byte) {
	for {
		select {
		case sw := <-r.Read:
			sshWrite <- sw
		}
	}
}

func (r *RelayClient) RelayClientReadSshWriteTcp(sshRead chan []byte) {
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

func (r *RelayClient) relayClientRead() {
	for {
		data := make([]byte, 1024)
		size, err := r.Conn.Read(data)
		if (err != nil) || (err == io.EOF) {
			if err != io.EOF {
				log.Error(err)
			}
			r.Conn.Close()
			return
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
				r.Conn.Close()
			}
		}
	}
}
