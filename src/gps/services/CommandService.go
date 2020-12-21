package services

import (
	"github.com/go-xorm/xorm"
	"goPanel/src/gps/models"
)

type CommandService struct {
	commandModel *models.CommandModel
}

func (s *CommandService) Add(db *xorm.Engine, data *models.CommandModel) (int64, error) {
	return s.commandModel.Add(db, data)
}

func (s *CommandService) IdByDetails(db *xorm.Engine, id int64) models.CommandModel {
	return s.commandModel.IdByDetails(db, id)
}

func (s *CommandService) Update(db *xorm.Engine, data models.CommandModel) (int64, error) {
	return s.commandModel.Update(db, data)
}
