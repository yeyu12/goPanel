package main

import (
	"github.com/gin-gonic/gin"
	core "goPanel/src/panel/core/database"
	"goPanel/src/panel/models"
	"goPanel/src/panel/router"
)

func main() {
	g := gin.Default()
	g = (new(router.Route)).Init(g)
	createTable()
	_ = g.Run(":10010")
}

// 创建表
func createTable() {
	core.CreateTables(
		new(models.MachineModel),
		new(models.MachineGroupModel),
	)
}
