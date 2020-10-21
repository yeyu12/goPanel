package services

import (
	"github.com/go-xorm/xorm"
	"goPanel/src/panel/models"
)

type MachineGroupService struct {
	machineGroupModel *models.MachineGroupModel
}

func (s *MachineGroupService) Add(db *xorm.Engine, data models.MachineGroupModel) (int64, error) {
	return s.machineGroupModel.Add(db, data)
}

func (s *MachineGroupService) Update(db *xorm.Engine, data models.MachineGroupModel) (int64, error) {
	return s.machineGroupModel.Update(db, data)
}

func (s *MachineGroupService) IdByDetails(db *xorm.Engine, id int64) models.MachineGroupModel {
	return s.machineGroupModel.IdByDetails(db, id)
}
