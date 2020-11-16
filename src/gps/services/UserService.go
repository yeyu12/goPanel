package services

import (
	"github.com/go-xorm/xorm"
	core "goPanel/src/core/database"
	"goPanel/src/gps/constants"
	"goPanel/src/gps/models"
	"time"
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

func (s *UserService) UpdateUser(db *xorm.Engine, data models.UserModel) (int64, error) {
	return s.userModel.UpdateUser(db, data)
}

// 是否登录，token是否有效
func (s *UserService) IsUserLogin(token string) (state bool, msg string, code int32) {
	if token == "" {
		return false, constants.PLEASE_LOG_IN_FIRST_MSG, constants.PLEASE_LOG_IN_FIRST
	} else {
		userData := s.TokenByData(core.Db, token)
		if userData.Id == 0 {
			return false, constants.TOKEN_INVALID_MSG, constants.TOKEN_INVALID
		} else if (time.Now().Unix() - userData.TokenExpirationTime.Unix()) > 86400 {
			return false, constants.TOKEN_BE_OVERDUE_MSG, constants.TOKEN_BE_OVERDUE
		}

		return true, constants.SUCCESS_MSG, constants.SUCCESS
	}
}
