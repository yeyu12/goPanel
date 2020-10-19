package models

import "time"

// 主机目录
type MachineGroupModel struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"create_time"`
	CreateUid  int64     `json:"create_uid"`
}
