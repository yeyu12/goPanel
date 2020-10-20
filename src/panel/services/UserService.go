package services

import (
	"github.com/go-xorm/xorm"
	"goPanel/src/panel/models"
)

type UserService struct {
	userModel *models.UserModel
}

func (s *UserService) UsernameAndPasswdByData(db *xorm.Engine, data map[string]string) models.UserModel {
	return s.userModel.UsernameAndPasswdByData(db, data)
}

func (s *UserService) UserAdd(db *xorm.Engine, data models.UserModel) (int64, error) {
	return s.userModel.UserAdd(db, data)
}

func (s *UserService) UsernameData(db *xorm.Engine, username string) models.UserModel {
	return s.userModel.UsernameData(db, username)
}

func (s *UserService) TokenByData(db *xorm.Engine, token string) models.UserModel {
	return s.userModel.TokenByData(db, token)
}
