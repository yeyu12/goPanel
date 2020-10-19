package main

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/config"
	core "goPanel/src/panel/core/database"
	core_log "goPanel/src/panel/core/log"
	"goPanel/src/panel/models"
	"goPanel/src/panel/router"
)

func main() {
	core_log.LogSetOutput(config.Conf.App.LogPath)
	createTable()

	g := gin.Default()
	g = (new(router.Route)).Init(g)
	_ = g.Run(":10010")
}

// 创建表
func createTable() {
	core.CreateTables(
		new(models.MachineModel),
		new(models.MachineGroupModel),
	)
}
