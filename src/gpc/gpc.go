package main

import (
	core_log "goPanel/src/core/log"
	"goPanel/src/gpc/config"
	"goPanel/src/gpc/service/socket"
	"time"
)

func main() {
	_, _ = time.LoadLocation("Asia/Shanghai")
	conf := config.Conf.App
	core_log.Initialization(
		conf.LogOutputType,
		conf.Debug,
		conf.LogLevel,
	)
	core_log.LogSetOutput(conf.LogPath, conf.LogOutputFlag)

	socket.ControlAddr = conf.ServerHost + ":" + conf.ServerPort
	socket.StartClientTcp()
}
