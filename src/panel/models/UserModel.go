package models

import (
	"github.com/go-xorm/xorm"
	"time"
)

// 用户
type UserModel struct {
	Id                  int64     `json:"id"`
	Username            string    `json:"username"`
	Passwd              string    `json:"passwd"`
	Token               string    `json:"token"`
	TokenExpirationTime time.Time `json:"token_expiration_time"` // token过期时间
	CreateTime          time.Time `json:"create_time"`
	UpdateTime          time.Time `json:"update_time"`
}

func (m *UserModel) UsernameAndPasswdByData(db *xorm.Engine, data map[string]string) UserModel {
	var user UserModel
	db.Where("username = ?", data["username"]).Where("passwd = ?", data["passwd"]).Get(&user)
	user.Passwd = ""

	return user
}

func (m *UserModel) UserAdd(db *xorm.Engine, data UserModel) (int64, error) {
	id, err := db.InsertOne(data)

	return id, err
}

func (m *UserModel) UsernameData(db *xorm.Engine, username string) UserModel {
	var user UserModel
	db.Where("username = ?", username).Get(&user)
	user.Passwd = ""

	return user
}

func (m *UserModel) TokenByData(db *xorm.Engine, username string) UserModel {
	var user UserModel
	db.Where("token = ?", username).Get(&user)
	user.Passwd = ""

	return user
}

func (m *UserModel) UpdateUser(db *xorm.Engine, data UserModel) (affected int64, err error) {
	affected, err = db.Where("id = ?", data.Id).Update(data)
	return
}
