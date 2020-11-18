package controltcp

import (
	log "github.com/sirupsen/logrus"
	"net"
)

func SettingInit(conn *net.TCPConn, message interface{}) {

}

func RegisterNode(conn *net.TCPConn, message interface{}) {
	log.Error(message)
}
