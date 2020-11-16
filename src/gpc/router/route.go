package router

import (
	"goPanel/src/gpc/controller"
	"net"
)

type handle func(conn *net.TCPConn, message interface{})

var Route = make(map[string]handle)

func init() {
	Route["connRelayByLocalSsh"] = controller.ConnRelayAndLocalSsh
}
