package controltcp

import (
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
	core "goPanel/src/core/database"
	"goPanel/src/gps/coer/socket"
	"goPanel/src/gps/services"
	"time"
	"unsafe"
)

func ExecCommandResult(cli unsafe.Pointer, message interface{}) {
	messBody, err := common.InterfaceByMapStr(message.(*socket.Message).Data)
	if err != nil {
		log.Error(err)
		return
	}

	commandService := services.CommandService{}
	oldCommandData := commandService.IdByDetails(core.Db, int64(messBody["id"].(float64)))
	if oldCommandData.Id == 0 {
		log.Errorf("执行的命令数据不存在！ %+v ", messBody)
		return
	}

	execTimeStr, _ := time.Parse(time.RFC3339, messBody["exec_time"].(string))
	oldCommandData.ExecTime = execTimeStr
	oldCommandData.IsExec = 1
	oldCommandData.ExecResult = messBody["exec_result"].(string)
	oldCommandData.HandleTime = int64(messBody["handle_time"].(float64))
	oldCommandData.UpdateTime = time.Now()

	_, err = commandService.Update(core.Db, oldCommandData)
	if err != nil {
		log.Errorf("执行命令，数据库状态更新失败！ %s ", err.Error())
		return
	}
}
