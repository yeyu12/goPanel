package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"goPanel/src/common"
	"goPanel/src/constants"
	"goPanel/src/gps/coer/socket"
	"goPanel/src/gps/validations"
	"io/ioutil"
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
			"id":          index.ClientId,
			"name":        index.Name,
			"system_type": index.SystemType,
		})
	}

	common.RetJson(g, constants.SUCCESS, constants.SUCCESS_MSG, ret)
	return
}

func (c *MachineController) Save(g *gin.Context) {
	inputData, _ := ioutil.ReadAll(g.Request.Body)
	var machineSaveComputerVail validations.MachineSaveComputer
	c.JsonPost(&machineSaveComputerVail, inputData)

	if err := c.Validations(machineSaveComputerVail); err != nil {
		common.RetJson(g, constants.ERROR_FAIL, err.Error(), "")
		return
	}

	cli := socket.ControlManager.FindClientIdByClientConn(machineSaveComputerVail.Id)
	if cli != nil {
		msg, _ := json.Marshal(socket.Message{
			Type:  0,
			Event: "settingClientInfo",
			Data:  machineSaveComputerVail,
			Code:  constants.SUCCESS,
		})

		cli.Write <- msg
	}

	common.RetJson(g, constants.SUCCESS, constants.SUCCESS_MSG, "")
	return
}

// 重启客户端主机
func (c *MachineController) Reboot(g *gin.Context) {
	clientId := g.Query("id")
	cli := socket.ControlManager.FindClientIdByClientConn(clientId)
	if cli != nil {
		msg, _ := json.Marshal(socket.Message{
			Type:  0,
			Event: "reboot",
			Data:  nil,
			Code:  constants.SUCCESS,
		})

		cli.Write <- msg
	}

	common.RetJson(g, constants.SUCCESS, constants.SUCCESS_MSG, "")
	return
}

// 重启客户单服务
func (c *MachineController) RestartService(g *gin.Context) {
	clientId := g.Query("id")
	cli := socket.ControlManager.FindClientIdByClientConn(clientId)
	if cli != nil {
		msg, _ := json.Marshal(socket.Message{
			Type:  0,
			Event: "restartService",
			Data:  nil,
			Code:  constants.SUCCESS,
		})

		cli.Write <- msg
	}

	// 停止监听指定中继端口
	for conn, _ := range socket.ServerWsManager.Clients {
		if conn.ClientId == clientId {
			socket.ServerWsManager.UnRegister <- conn
		}
	}

	common.RetJson(g, constants.SUCCESS, constants.SUCCESS_MSG, "")
	return
}
