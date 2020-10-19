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

func (m *UserModel) UsernameAndPasswdByData(db *xorm.Engine, data map[string]string) *UserModel {
	var user *UserModel
	_ = db.Where("username = ?", data["username"]).Where("passwd = ?", data["passwd"]).Find(user)

	return user
}
