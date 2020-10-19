package models

import "time"

// 主机
type MachineModel struct {
	Id             int64     `json:"id"`
	MachineGroupId int64     `json:"machine_group_id"`
	Alias          string    `json:"alias"`       //别名
	Host           string    `json:"host"`        // 地址
	User           string    `json:"user"`        // 用户名
	Port           int       `json:"port"`        // 端口
	CreateTime     time.Time `json:"create_time"` // 创建时间
	LoginNum       int64     `json:"login_num"`   // 登录次数
	UpdateTime     time.Time `json:"update_time"` // 更新次数
	CreateUid      int64     `json:"create_uid"`  // 创建人
}
