package models

import (
	"github.com/go-xorm/xorm"
	"time"
)

// 主机
type MachineModel struct {
	Id             int64     `json:"id"`
	MachineGroupId int64     `json:"machine_group_id"`
	Name           string    `json:"name"`        //别名
	Host           string    `json:"host"`        // 地址
	User           string    `json:"user"`        // 用户名
	Port           int       `json:"port"`        // 端口
	CreateTime     time.Time `json:"create_time"` // 创建时间
	LoginNum       int64     `json:"login_num"`   // 登录次数
	UpdateTime     time.Time `json:"update_time"` // 更新次数
	CreateUid      int64     `json:"create_uid"`  // 创建人
}

func (m *MachineModel) Get(db *xorm.Engine, where map[string]interface{}) *[]MachineModel {
	var data []MachineModel
	dbs := db.Asc("id")
	if where["machine_group_id"] != nil {
		dbs = db.Where("machine_group_id = ?", where["machine_group_id"])
	}
	dbs.Find(&data)

	return &data
}

func (m *MachineModel) Add(db *xorm.Engine, data MachineModel) (int64, error) {
	id, err := db.InsertOne(data)

	return id, err
}

func (m *MachineModel) Update(db *xorm.Engine, data MachineModel) (affected int64, err error) {
	affected, err = db.Where("id = ?", data.Id).Update(data)
	return
}

func (m *MachineModel) IdByDetails(db *xorm.Engine, id int64) MachineModel {
	var machineGroupData MachineModel
	db.Where("id = ?", id).Get(&machineGroupData)

	return machineGroupData
}
