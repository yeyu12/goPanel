package main

import (
	core_log "goPanel/src/core/log"
	"goPanel/src/gpc/config"
	"goPanel/src/gpc/service/socket"
	"time"
)

func main() {
	_, _ = time.LoadLocation("Asia/Shanghai")
	core_log.Initialization(
		config.Conf.App.LogOutputType,
		config.Conf.App.Debug,
		config.Conf.App.LogLevel,
	)
	core_log.LogSetOutput(config.Conf.App.LogPath, config.Conf.App.LogOutputFlag)

	socket.ControlAddr = config.Conf.App.ServerHost + ":" + config.Conf.App.ServerPort
	socket.StartClientTcp()
}
