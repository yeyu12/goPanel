package socket

import (
	"fmt"
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
		fmt.Println("中继端没有端口可用！")
		return -1
	}

	return ControlManager.relayStartPort
}

func CreateRelayConn(port int) {
	relayListen, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("创建中继端失败", err)
	}

	go func() {
		for {
			relayConn, err := relayListen.Accept()
			if err != nil {
				fmt.Println("连接中继端失败!", err)
			}

			for {
				data := make([]byte, 10240)
				size, err := relayConn.Read(data)
				if err != nil {
					fmt.Println("消息读失败！", err)
				}
				data = data[:size]
				fmt.Println("来自用户端：", relayConn.RemoteAddr())
				fmt.Println("中继端的输出：", string(data))
			}
		}
	}()
}
