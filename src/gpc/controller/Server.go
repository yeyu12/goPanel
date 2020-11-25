package controller

import (
	"encoding/json"
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
	err := relayClient.RelayConn(relayAddr)
	if err != nil {
		log.Error(err)
		return
	}
	/*defer func() {
		if err := relayClient.Conn.Close(); err != nil {
			log.Error(err)
			return
		}
	}()*/

	tcpSsh := ssh.NewTcpSsh()
	sh, sshChannel, err := tcpSsh.SshConn("127.0.0.1", "fengxiao", "ZpB123", 22, uint32(data["cols"].(float64)), uint32(data["rows"].(float64)))
	if err != nil {
		m := service.Message{
			Event: constants.WS_EVENT_ERR,
			Data:  constants.SSH_CONNECTION_FAILED_MSG,
			Code:  constants.SSH_CONNECTION_FAILED,
		}

		resJson, _ := json.Marshal(m)
		_, err = conn.Write(resJson)
		if err != nil {
			log.Error(err)
		}

		return
	}
	/*defer func() {
		if err := sshChannel.Close(); err != nil {
			log.Error(err)
			return
		}
	}()*/

	// tcp和ssh交换数据
	go relayClient.RelayClientReadTcpWriteSsh(tcpSsh.SshWrite)
	go relayClient.RelayClientReadSshWriteTcp(tcpSsh.SshRead)

	go sh.Read(sshChannel, tcpSsh.SshRead)
	go sh.Write(sshChannel, tcpSsh.SshWrite)
}
