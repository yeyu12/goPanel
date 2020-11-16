package main

import (
	core_log "goPanel/src/core/log"
	"goPanel/src/gpc/config"
	"goPanel/src/gpc/service/socket"
	"time"
)

func main() {
	_, _ = time.LoadLocation("Asia/Shanghai")
	core_log.LogSetOutput(config.Conf.App.LogPath)

	socket.StartClientTcp(config.Conf.App.ServerHost + ":" + config.Conf.App.ServerPort)
}
