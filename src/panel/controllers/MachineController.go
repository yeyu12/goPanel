package controllers

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/common"
	"goPanel/src/panel/constants"
	core "goPanel/src/panel/core/database"
	"goPanel/src/panel/models"
	"goPanel/src/panel/services"
	"goPanel/src/panel/validations"
	"io/ioutil"
	"time"
)

type MachineController struct {
	BaseController
	machineService      *services.MachineService
	machineGroupService *services.MachineGroupService
}

func NewMachineController() *MachineController {
	return &MachineController{
		machineService:      new(services.MachineService),
		machineGroupService: new(services.MachineGroupService),
	}
}

func (c *MachineController) List(g *gin.Context) {
	common.RetJson(g, 200, "成功", "")
}

func (c *MachineController) Add(g *gin.Context) {
	inputData, _ := ioutil.ReadAll(g.Request.Body)
	var addVail validations.MachineAdd
	c.JsonPost(&addVail, inputData)

	if err := c.Validations(addVail); err != nil {
		common.RetJson(g, constants.MISSING_PARAMETER_FAIL, err.Error(), "")
		return
	}

	var (
		code int32
		msg  string
		data interface{}
	)
	switch addVail.Flag {
	case 1:
		code, msg, data = c.addDir(g, inputData)
		if code != constants.SUCCESS {
			common.RetJson(g, code, msg, data)
			return
		}

		break
	case 2:
		code, msg, data = c.addComputer(g, inputData)
		if code != constants.SUCCESS {
			common.RetJson(g, code, msg, data)
			return
		}

		break
	}

	common.RetJson(g, code, msg, data)
	return
}

func (c *MachineController) Edit(g *gin.Context) {
	common.RetJson(g, 200, "成功", "")
}

func (c *MachineController) Del(g *gin.Context) {
	common.RetJson(g, 200, "成功", "")
}

// 添加目录
func (c *MachineController) addDir(g *gin.Context, inputData []byte) (int32, string, interface{}) {
	var addDirVail validations.MachineAddDir
	c.JsonPost(&addDirVail, inputData)

	if err := c.Validations(addDirVail); err != nil {
		return constants.MISSING_PARAMETER_FAIL, err.Error(), ""
	}

	userinfo := c.GetUserInfo(g)
	var addDirData models.MachineGroupModel
	c.JsonPost(&addDirData, inputData)
	addDirData.CreateTime = time.Now()
	addDirData.CreateUid = userinfo.Id

	id, err := c.machineGroupService.Add(core.Db, addDirData)
	if err != nil {
		return constants.ERROR_FAIL, constants.ERROR_FAIL_MSG, ""
	}

	data := c.machineGroupService.IdByDetails(core.Db, id)

	return constants.SUCCESS, constants.SUCCESS_MSG, data
}

// 添加主机
func (c *MachineController) addComputer(g *gin.Context, inputData []byte) (int32, string, interface{}) {
	var addComputerVail validations.MachineAddComputer
	c.JsonPost(&addComputerVail, inputData)

	if err := c.Validations(addComputerVail); err != nil {
		return constants.MISSING_PARAMETER_FAIL, err.Error(), ""
	}

	userinfo := c.GetUserInfo(g)
	var addComputerData models.MachineModel
	c.JsonPost(&addComputerData, inputData)
	addComputerData.CreateTime = time.Now()
	addComputerData.CreateUid = userinfo.Id
	addComputerData.LoginNum = 0

	id, err := c.machineService.Add(core.Db, addComputerData)
	if err != nil {
		return constants.ERROR_FAIL, constants.ERROR_FAIL_MSG, ""
	}

	data := c.machineService.IdByDetails(core.Db, id)

	return constants.SUCCESS, constants.SUCCESS_MSG, data
}
