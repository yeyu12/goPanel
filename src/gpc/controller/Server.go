package controller

import (
	log "github.com/sirupsen/logrus"
	"net"
)

// 连接中继端，和本地ssh连接
func ConnRelayAndLocalSsh(conn *net.TCPConn, message interface{}) {
	log.Info("进入执行器")
	log.Error(message)
}

//
