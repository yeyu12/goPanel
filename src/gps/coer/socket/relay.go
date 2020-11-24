package socket

import (
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
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

func CreateRelayConn(port int) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	relayListen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			relayConn, err := relayListen.Accept()
			if err != nil {
				log.Info("连接中继端失败!", err)
			}

			for {
				data := make([]byte, 10240)
				size, err := relayConn.Read(data)
				if err != nil {
					log.Error("消息读失败！", err)
				}
				data = data[:size]
				log.Info("来自用户端：", relayConn.RemoteAddr())
				log.Info("中继端的输出：", string(data))
			}
		}
	}()

	return relayListen, nil
}
