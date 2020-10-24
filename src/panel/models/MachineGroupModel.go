package models

import (
	"github.com/go-xorm/xorm"
	"time"
)

// 主机目录
type MachineGroupModel struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	IsDir      int       `json:"is_dir" xorm:"default(1)"`
	CreateTime time.Time `json:"create_time"`
	CreateUid  int64     `json:"create_uid"`
	UpdateTime time.Time `json:"update_time"`
}

func (m *MachineGroupModel) Get(db *xorm.Engine) *[]MachineGroupModel {
	var data []MachineGroupModel
	db.Asc("id").Find(&data)

	return &data
}

func (m *MachineGroupModel) Add(db *xorm.Engine, data MachineGroupModel) (int64, error) {
	data.IsDir = 1
	id, err := db.InsertOne(data)

	return id, err
}

func (m *MachineGroupModel) Update(db *xorm.Engine, data MachineGroupModel) (affected int64, err error) {
	affected, err = db.Where("id = ?", data.Id).Update(data)
	return
}

func (m *MachineGroupModel) Del(db *xorm.Engine, id int64) (affected int64, err error) {
	affected, err = db.Id(id).Delete(new(MachineGroupModel))
	return
}

func (m *MachineGroupModel) IdByDetails(db *xorm.Engine, id int64) MachineGroupModel {
	var machineGroupData MachineGroupModel
	db.Where("id = ?", id).Get(&machineGroupData)

	return machineGroupData
}
