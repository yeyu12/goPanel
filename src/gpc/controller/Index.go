package controller

import (
	log "github.com/sirupsen/logrus"
	"net"
)

func Heartbeat(conn *net.TCPConn, message interface{}) {
	log.Info("心跳")
}
