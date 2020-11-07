package models

import (
	"github.com/go-xorm/xorm"
	"time"
)

// 主机
type MachineModel struct {
	Id             int64     `json:"id"`
	MachineGroupId int64     `json:"machine_group_id" xorm:"index default(0)"`
	Name           string    `json:"name" xorm:"varchar(50)"`     //别名
	Host           string    `json:"host" xorm:"char(16) index"`  // 地址
	User           string    `json:"user" xorm:"default('root')"` // 用户名
	Port           int       `json:"port" xorm:"default(22)"`     // 端口
	CreateTime     time.Time `json:"create_time"`                 // 创建时间
	LoginNum       int64     `json:"login_num" xorm:"default(0)"` // 登录次数
	UpdateTime     time.Time `json:"update_time"`                 // 更新次数
	CreateUid      int64     `json:"create_uid"`                  // 创建人
}

func (m *MachineModel) Get(db *xorm.Engine, where map[string]interface{}) *[]MachineModel {
	var data []MachineModel
	dbs := db.Asc("id")
	if where["machine_group_id"] != nil {
		dbs = db.Where("machine_group_id = ?", where["machine_group_id"])
	} else {
		dbs = db.Where("machine_group_id = ?", 0)
	}
	dbs.Find(&data)

	return &data
}

func (m *MachineModel) GetAll(db *xorm.Engine) *[]MachineModel {
	var data []MachineModel
	db.Find(&data)

	return &data
}

func (m *MachineModel) Add(db *xorm.Engine, data *MachineModel) (int64, error) {
	id, err := db.InsertOne(data)

	return id, err
}

func (m *MachineModel) Update(db *xorm.Engine, data MachineModel) (affected int64, err error) {
	affected, err = db.Where("id = ?", data.Id).Update(data)
	return
}

func (m *MachineModel) Del(db *xorm.Engine, id int64) (affected int64, err error) {
	affected, err = db.Id(id).Delete(new(MachineModel))
	return
}

func (m *MachineModel) IdByDetails(db *xorm.Engine, id int64) MachineModel {
	var machineGroupData MachineModel
	db.Where("id = ?", id).Get(&machineGroupData)

	return machineGroupData
}
