package controllers

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/common"
	"goPanel/src/gps/coer/socket"
	"goPanel/src/gps/constants"
)

type MachineController struct {
	BaseController
}

func NewMachineController() *MachineController {
	return &MachineController{}
}

func (c *MachineController) List(g *gin.Context) {
	var ret []map[string]interface{}
	for index, _ := range socket.ControlManager.Clients {
		ret = append(ret, map[string]interface{}{
			"id":   index.Uuid,
			"name": index.Name,
		})
	}

	common.RetJson(g, constants.SUCCESS, constants.SUCCESS_MSG, ret)
	return
}
