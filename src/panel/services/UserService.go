package services

import (
	"github.com/go-xorm/xorm"
	"goPanel/src/panel/models"
)

type UserService struct {
	userModel *models.UserModel
}

func (s *UserService) UsernameAndPasswdByData(db *xorm.Engine, data map[string]string) *models.UserModel {
	return s.userModel.UsernameAndPasswdByData(db, data)
}
