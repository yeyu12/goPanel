package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

const (
	CREATE_DIR      = 1
	CREATE_COMPUTER = 2
)

func NewMachineController() *MachineController {
	return &MachineController{
		machineService:      new(services.MachineService),
		machineGroupService: new(services.MachineGroupService),
	}
}

func (c *MachineController) List(g *gin.Context) {
	groupData := c.machineGroupService.Get(core.Db)
	var retData []map[string]interface{}
	groupJson, _ := json.Marshal(groupData)
	_ = json.Unmarshal(groupJson, &retData)

	for index, item := range retData {
		where := map[string]interface{}{
			"machine_group_id": item["id"],
		}

		retData[index]["children"] = c.machineService.Get(core.Db, where)
	}

	machineData := c.machineService.Get(core.Db, map[string]interface{}{})
	var machineMap []map[string]interface{}
	machineJson, _ := json.Marshal(machineData)
	_ = json.Unmarshal(machineJson, &machineMap)

	retData = append(retData, machineMap...)

	common.RetJson(g, constants.SUCCESS, constants.SUCCESS_MSG, retData)
	return
}

func (c *MachineController) Save(g *gin.Context) {
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
	case CREATE_DIR:
		code, msg, data = c.saveDir(g, inputData)
		if code != constants.SUCCESS {
			common.RetJson(g, code, msg, data)
			return
		}

		break
	case CREATE_COMPUTER:
		code, msg, data = c.saveComputer(g, inputData)
		if code != constants.SUCCESS {
			common.RetJson(g, code, msg, data)
			return
		}

		break
	}

	common.RetJson(g, code, msg, data)
	return
}

func (c *MachineController) Del(g *gin.Context) {
	inputData, _ := ioutil.ReadAll(g.Request.Body)
	var delVail validations.MachineDel
	c.JsonPost(&delVail, inputData)

	if err := c.Validations(delVail); err != nil {
		common.RetJson(g, constants.MISSING_PARAMETER_FAIL, err.Error(), "")
		return
	}

	log.Error(delVail)

	switch delVail.Flag {
	case CREATE_DIR:
		_, err := c.machineGroupService.Del(core.Db, delVail.Id)
		if err != nil {
			common.RetJson(g, constants.ERROR_FAIL, constants.ERROR_FAIL_MSG, "")
			return
		}

		break
	case CREATE_COMPUTER:
		_, err := c.machineService.Del(core.Db, delVail.Id)
		if err != nil {
			common.RetJson(g, constants.ERROR_FAIL, constants.ERROR_FAIL_MSG, "")
			return
		}

		break
	}

	common.RetJson(g, constants.SUCCESS, constants.SUCCESS_MSG, "")
	return
}

// 添加目录
func (c *MachineController) saveDir(g *gin.Context, inputData []byte) (int32, string, interface{}) {
	var addDirVail validations.MachineAddDir
	c.JsonPost(&addDirVail, inputData)

	if err := c.Validations(addDirVail); err != nil {
		return constants.MISSING_PARAMETER_FAIL, err.Error(), ""
	}

	userinfo := c.GetUserInfo(g)
	var addDirData models.MachineGroupModel
	c.JsonPost(&addDirData, inputData)
	var retData interface{}

	if addDirData.Id == 0 {
		addDirData.CreateTime = time.Now()
		addDirData.CreateUid = userinfo.Id

		id, err := c.machineGroupService.Add(core.Db, addDirData)
		if err != nil {
			return constants.ERROR_FAIL, constants.ERROR_FAIL_MSG, ""
		}

		addDirData.Id = id
	} else {
		addDirData.UpdateTime = time.Now()
		_, err := c.machineGroupService.Update(core.Db, addDirData)
		if err != nil {
			return constants.ERROR_FAIL, constants.ERROR_FAIL_MSG, ""
		}
	}

	retData = c.machineGroupService.IdByDetails(core.Db, addDirData.Id)

	return constants.SUCCESS, constants.SUCCESS_MSG, retData
}

// 添加主机
func (c *MachineController) saveComputer(g *gin.Context, inputData []byte) (int32, string, interface{}) {
	var addComputerVail validations.MachineAddComputer
	c.JsonPost(&addComputerVail, inputData)

	if err := c.Validations(addComputerVail); err != nil {
		return constants.MISSING_PARAMETER_FAIL, err.Error(), ""
	}

	userinfo := c.GetUserInfo(g)
	var addComputerData models.MachineModel
	c.JsonPost(&addComputerData, inputData)

	if addComputerData.Id == 0 {
		addComputerData.CreateTime = time.Now()
		addComputerData.CreateUid = userinfo.Id
		addComputerData.LoginNum = 0

		id, err := c.machineService.Add(core.Db, addComputerData)
		if err != nil {
			return constants.ERROR_FAIL, constants.ERROR_FAIL_MSG, ""
		}

		addComputerData.Id = id
	} else {
		addComputerData.UpdateTime = time.Now()
		_, err := c.machineService.Update(core.Db, addComputerData)
		if err != nil {
			return constants.ERROR_FAIL, constants.ERROR_FAIL_MSG, ""
		}
	}

	data := c.machineService.IdByDetails(core.Db, addComputerData.Id)
	dataMap := common.StructToJson(data)
	rsaPublicKey, err := ioutil.ReadFile(common.GetRsaFilePath() + "public.pem")
	if err != nil {
		log.Error(err)
	}
	encodePasswd, err := common.RsaEncrypt([]byte(addComputerVail.Passwd), rsaPublicKey)
	if err != nil {
		log.Error(err)
	}
	dataMap["passwd"] = encodePasswd

	return constants.SUCCESS, constants.SUCCESS_MSG, dataMap
}
