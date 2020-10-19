package controllers

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/common"
	"goPanel/src/panel/validations"
	"io/ioutil"
)

type MachineController struct {
	BaseController
}

func NewMachineController() *MachineController {
	return &MachineController{}
}

func (c *MachineController) MachineList(g *gin.Context) {
	common.RetJson(g, 200, "成功", "")
}

func (c *MachineController) MachineAdd(g *gin.Context) {
	inputData, _ := ioutil.ReadAll(g.Request.Body)
	var userVail validations.MachineAdd
	c.JsonPost(&userVail, inputData)

	if err := c.Validations(g, userVail); err != nil {
		common.RetJson(g, 4000, err.Error(), "")
		return
	}

	common.RetJson(g, 200, "成功", "")
}

func (c *MachineController) MachineEdit(g *gin.Context) {
	common.RetJson(g, 200, "成功", "")
}

func (c *MachineController) MachineDel(g *gin.Context) {
	common.RetJson(g, 200, "成功", "")
}
