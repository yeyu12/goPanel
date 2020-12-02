package controller

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net"
)

func Heartbeat(ctx context.Context, conn *net.TCPConn, message interface{}) {
	log.Info("心跳")
}
