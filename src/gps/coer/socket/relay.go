package socket

import (
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
	"io"
	"net"
	"strconv"
)

type Relay struct {
	Conn *net.TCPConn
}

func (r *Relay) RelayPort() int {
	if len(ControlManager.recoveryPort) > 0 {
		retPort := ControlManager.recoveryPort[0]
		ControlManager.recoveryPort = ControlManager.recoveryPort[1:]

		return retPort
	}

	// 中继端口使用监测
	ControlManager.relayStartPort = common.RetRelayPort(ControlManager.relayStartPort)
	if ControlManager.relayStartPort == -1 {
		log.Error("中继端没有端口可用！")
		return -1
	}

	return ControlManager.relayStartPort
}

func (r *Relay) CreateRelayConn(port int, wsWrite, wsRead chan []byte, relayConnCh chan *net.TCPConn) (*net.TCPListener, error) {
	relayAddr := "0.0.0.0:" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", relayAddr)
	if err != nil {
		return nil, err
	}

	relayListen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	log.Info("中继端启动：", relayAddr)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
				relayListen.Close()
			}
		}()

		for {
			relayConn, err := relayListen.Accept()
			if err != nil {
				log.Info(err)
				break
			}

			r.Conn = relayConn.(*net.TCPConn)
			relayConnCh <- r.Conn

			go r.read(wsWrite)
			go r.write(wsRead)
		}
	}()

	return relayListen, nil
}

func (r *Relay) read(wsWrite chan []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			r.Conn.Close()
		}
	}()

	for {
		data := make([]byte, 10240)
		size, err := r.Conn.Read(data)
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Info("消息读失败！", err)
			break
		}
		data = data[:size]

		log.Info("来自用户端的消息：", r.Conn.RemoteAddr())
		log.Info("消息内容：", string(data))

		// 中继端的输出发送到ws中
		wsWrite <- data
	}
}

func (r *Relay) write(wsRead chan []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			r.Conn.Close()
		}
	}()

	for {
		select {
		case read := <-wsRead:
			_, err := r.Conn.Write(read)
			if err != nil {
				log.Info(err)
				return
			}
		}
	}
}
