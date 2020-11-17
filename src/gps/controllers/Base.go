package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"goPanel/src/gps/models"
	"goPanel/src/gps/validations"
)

type BaseController struct {
}

func (c *BaseController) JsonPost(dataStruct interface{}, data []byte) interface{} {
	_ = json.Unmarshal(data, &dataStruct)
	return dataStruct
}

func (c *BaseController) Validations(vali interface{}) error {
	if err := validations.Validate.Struct(vali); err != nil {
		return validations.Translate(err.(validator.ValidationErrors))
	}

	return nil
}

func (c *BaseController) GetUserInfo(g *gin.Context) *models.UserModel {
	userinfo, exists := g.Get("userinfo")
	if !exists {
		return nil
	}

	return userinfo.(*models.UserModel)
}
