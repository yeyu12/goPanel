package controller

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"goPanel/src/gpc/service"
	"net"
)

func HandleCommand(ctx context.Context, conn *net.TCPConn, message interface{}) {
	dataMap := message.(service.Message).Data.(map[string]interface{})
	dataJson, _ := json.Marshal(dataMap)
	var command service.CommandService
	_ = json.Unmarshal(dataJson, &command)

	var tempWaiteExecCommand []service.CommandService
	if len(service.WaitExecCommandData) == 0 {
		tempWaiteExecCommand = append(tempWaiteExecCommand, command)
	} else {
		// 按时间排序
		for index, item := range service.WaitExecCommandData {
			if item.PlanExecTime.Unix() > command.PlanExecTime.Unix() {
				endWaiteExecCommand := make([]service.CommandService, 1)
				size := copy(endWaiteExecCommand, service.WaitExecCommandData[index:])
				log.Infof("复制切片 %d 个，切片值 %+v", size, endWaiteExecCommand)
				tempWaiteExecCommand = append(service.WaitExecCommandData[:index], command)
				tempWaiteExecCommand = append(tempWaiteExecCommand, endWaiteExecCommand...)
				goto JAMP
			}
		}
		tempWaiteExecCommand = append(service.WaitExecCommandData, command)

	JAMP:
	}

	service.WaitExecCommandData = tempWaiteExecCommand
}
