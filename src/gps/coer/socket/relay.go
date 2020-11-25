package socket

import (
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
	"io"
	"net"
	"strconv"
)

func RelayPort() int {
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

func CreateRelayConn(port int, wsWrite chan []byte) (*net.TCPListener, error) {
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
			}
		}()

		for {
			relayConn, err := relayListen.Accept()
			if err != nil {
				log.Info(err)
				return
			}

			go func() {
				for {
					data := make([]byte, 10240)
					size, err := relayConn.Read(data)
					if err == io.EOF {
						return
					}
					if err != nil {
						log.Error("消息读失败！", err)
					}
					data = data[:size]

					log.Info("来自用户端的消息：", relayConn.RemoteAddr())
					//log.Info("消息内容：", string(data))

					// 中继端的输出发送到ws中
					wsWrite <- data
				}
			}()
		}
	}()

	return relayListen, nil
}
