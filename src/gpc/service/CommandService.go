package service

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"goPanel/src/constants"
	"os/exec"
	"time"
)

var WaitExecCommandData []CommandService // 待执行命令

// 获取要执行的命令
func GetExecCommand() {
	for {
		localTime := time.Now().Unix()

		for _, item := range WaitExecCommandData {
			if !item.IsLock && localTime >= item.PlanExecTime.Unix() {
				go execCommand(&item)
			} else {
				break
			}
		}

		time.Sleep(time.Second * 1)
	}
}

// 执行命令
func execCommand(data *CommandService) {
	// 命令执行完成后，上传控制端，然后删除待执行的列表
	data.IsLock = true
	startTime := time.Now()
	cmd := exec.Command("/bin/bash", "-c", data.Command)
	output, err := cmd.Output()
	endTime := time.Now().UnixNano() / 1e6

	if err != nil {
		data.ExecResult = err.Error()
	} else {
		data.ExecResult = string(output)
	}

	data.ExecTime = startTime
	data.HandleTime = endTime - (startTime.UnixNano() / 1e6)
	data.IsExec = 1

	// 删除本地存储，向服务端推送
	// 找到下标
	index := 0
	for key, val := range WaitExecCommandData {
		if val.Id == data.Id {
			index = key
			break
		}
	}
	// 删除
	WaitExecCommandData = append(WaitExecCommandData[:index], WaitExecCommandData[index+1:]...)

	res := Message{
		Event: "execCommandResult",
		Data:  data,
		Code:  constants.SUCCESS,
	}

	resJson, _ := json.Marshal(res)
	err = NewTcpService(data.Conn).Send(resJson)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(data)
}
