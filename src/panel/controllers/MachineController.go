package controllers

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/common"
	"goPanel/src/panel/validations"
	"io/ioutil"
)

func MachineList(g *gin.Context) {
	common.RetJson(g, 200, "成功", "")
}

func MachineAdd(g *gin.Context) {
	inputData, _ := ioutil.ReadAll(g.Request.Body)
	var userVail validations.MachineAdd
	JsonPost(&userVail, inputData)

	if err := Validations(g, userVail); err != nil {
		common.RetJson(g, 4000, err.Error(), "")
		return
	}

	common.RetJson(g, 200, "成功", "")
}

func MachineEdit(g *gin.Context) {
	common.RetJson(g, 200, "成功", "")
}

func MachineDel(g *gin.Context) {
	common.RetJson(g, 200, "成功", "")
}
