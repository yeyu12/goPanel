package router

import (
	"goPanel/src/gps/controllers/controltcp"
	"net"
)

type handle func(conn *net.TCPConn, message interface{})

var ControlTcpRoute = make(map[string]handle)

func init() {
	ControlTcpRoute["local_setting"] = controltcp.SettingInit
	ControlTcpRoute["local_register"] = controltcp.RegisterNode
}
