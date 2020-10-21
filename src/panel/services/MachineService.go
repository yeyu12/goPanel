package services

import (
	"github.com/go-xorm/xorm"
	"goPanel/src/panel/models"
)

type MachineService struct {
	machineModel *models.MachineModel
}

func (s *MachineService) Add(db *xorm.Engine, data models.MachineModel) (int64, error) {
	return s.machineModel.Add(db, data)
}

func (s *MachineService) Update(db *xorm.Engine, data models.MachineModel) (int64, error) {
	return s.machineModel.Update(db, data)
}

func (s *MachineService) IdByDetails(db *xorm.Engine, id int64) models.MachineModel {
	return s.machineModel.IdByDetails(db, id)
}
