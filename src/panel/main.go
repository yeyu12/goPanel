package main

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/common"
	"goPanel/src/panel/config"
	core "goPanel/src/panel/core/database"
	core_log "goPanel/src/panel/core/log"
	"goPanel/src/panel/models"
	"goPanel/src/panel/router"
	"goPanel/src/panel/services/websocket"
	"time"
)

func main() {
	// 服务启动前，初始化操作
	_, _ = time.LoadLocation("Asia/Shanghai")
	core_log.LogSetOutput(config.Conf.App.LogPath)
	createTable()
	go websocket.WsManager.Start()
	common.GenRsaKey(common.GetRsaFilePath(), 2048)

	g := gin.Default()
	g = (new(router.Route)).Init(g)
	_ = g.Run(":10010")
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
