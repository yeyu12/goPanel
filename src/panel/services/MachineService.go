package services

import (
	"github.com/go-xorm/xorm"
	"goPanel/src/panel/models"
)

type MachineService struct {
	machineModel *models.MachineModel
}

func (s *MachineService) Get(db *xorm.Engine, where map[string]interface{}) *[]models.MachineModel {
	return s.machineModel.Get(db, where)
}

func (s *MachineService) GetAll(db *xorm.Engine) *[]models.MachineModel {
	return s.machineModel.GetAll(db)
}

func (s *MachineService) Add(db *xorm.Engine, data *models.MachineModel) (int64, error) {
	return s.machineModel.Add(db, data)
}

func (s *MachineService) Update(db *xorm.Engine, data models.MachineModel) (int64, error) {
	return s.machineModel.Update(db, data)
}

func (s *MachineService) Del(db *xorm.Engine, id int64) (int64, error) {
	return s.machineModel.Del(db, id)
}

func (s *MachineService) IdByDetails(db *xorm.Engine, id int64) models.MachineModel {
	return s.machineModel.IdByDetails(db, id)
}
