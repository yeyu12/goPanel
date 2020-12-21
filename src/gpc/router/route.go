package router

import (
	"context"
	"goPanel/src/gpc/controller"
	"net"
)

type handle func(context.Context, *net.TCPConn, interface{})

var Route = make(map[string]handle)

func init() {
	Route["sshConnectRelay"] = controller.SshConnectRelay
	Route["settingClientInfo"] = controller.SettingClientInfo
	Route["heartbeat"] = controller.Heartbeat
	Route["reboot"] = controller.Reboot
	Route["restartService"] = controller.RestartService
	Route["handleCommand"] = controller.HandleCommand
}
