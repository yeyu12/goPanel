package main

import (
	"github.com/gin-gonic/gin"
	core "goPanel/src/core/database"
	core_log "goPanel/src/core/log"
	"goPanel/src/gps/coer/socket"
	"goPanel/src/gps/config"
	"goPanel/src/gps/models"
	"goPanel/src/gps/router"
	"strconv"
	"time"
)

func main() {
	// 服务启动前，初始化操作
	_, _ = time.LoadLocation("Asia/Shanghai")
	conf := config.Conf.App
	core_log.Initialization(
		conf.LogOutputType,
		conf.Debug,
		conf.LogLevel,
	)
	core_log.LogSetOutput(conf.LogPath, conf.LogOutputFlag)
	createTable()
	go socket.ServerWsManager.Start()
	go socket.ControlManager.Start()
	//common.GenRsaKey(common.GetRsaFilePath(), 2048)

	g := gin.Default()
	g = (new(router.Route)).Init(g)
	_ = g.Run(":" + strconv.Itoa(conf.HttpPort))
}

// 创建表
func createTable() {
	core.CreateTables(
		new(models.UserModel),
		new(models.MachineModel),
		new(models.MachineGroupModel),
		new(models.CommandModel),
	)
}
