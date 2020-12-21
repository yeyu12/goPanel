package service

import (
	"time"
)

var WaitExecCommandData []CommandService // 待执行命令

// 获取要执行的命令
func GetExecCommand() {
	for {
		localTime := time.Now().Unix()

		for _, item := range WaitExecCommandData {
			if localTime >= item.PlanExecTime.Unix() {
				//log.Info("等待执行的命令：", item.Command)
				go execCommand(item)
			} else {
				break
			}
		}

		time.Sleep(time.Second * 1)
	}
}

// 执行命令
func execCommand(data CommandService) {
	// 命令执行完成后，上传控制端，然后删除待执行的列表
}
