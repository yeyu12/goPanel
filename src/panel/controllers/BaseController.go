package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"goPanel/src/panel/validations"
)

type BaseController struct {
}

func (c *BaseController) JsonPost(dataStruct interface{}, data []byte) interface{} {
	_ = json.Unmarshal(data, &dataStruct)
	return dataStruct
}

func (c *BaseController) Validations(g *gin.Context, vali interface{}) error {
	if err := validations.Validate.Struct(vali); err != nil {
		return validations.Translate(err.(validator.ValidationErrors))
	}

	return nil
}

func (c *BaseController) Panic(err error) {
	if err != nil {

	}
}
