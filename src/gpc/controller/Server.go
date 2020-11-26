package controller

import (
	log "github.com/sirupsen/logrus"
	"goPanel/src/constants"
	"goPanel/src/gpc/config"
	"goPanel/src/gpc/service"
	"goPanel/src/gpc/service/ssh"
	"net"
	"strconv"
)

// 连接中继端，和本地ssh连接
func SshConnectRelay(conn *net.TCPConn, message interface{}) {
	log.Info("进入执行器")

	data := message.(service.Message).Data.(map[string]interface{})
	relayClient := ssh.NewRelayClient()
	relayAddr := config.Conf.App.ServerHost + ":" + strconv.Itoa(int(data["port"].(float64)))
	log.Info("连接中继端，relayAddr:", relayAddr)
	err := relayClient.RelayConn(relayAddr, constants.CLIENT_SHELL_TYPE, uint32(data["cols"].(float64)), uint32(data["rows"].(float64)))
	if err != nil {
		log.Error(err)
		return
	}
}
