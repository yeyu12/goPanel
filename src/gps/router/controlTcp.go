package router

import (
	"goPanel/src/gps/coer/router"
	"goPanel/src/gps/controllers/controltcp"
)

func init() {
	router.AddRoute("local_setting", controltcp.SettingInit)
	router.AddRoute("local_register", controltcp.RegisterNode)
	router.AddRoute("heartbeat", controltcp.Heartbeat)
	router.AddRoute("execCommandResult", controltcp.ExecCommandResult)
}
