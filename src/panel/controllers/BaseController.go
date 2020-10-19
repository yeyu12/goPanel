package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"goPanel/src/panel/validations"
)

func JsonPost(dataStruct interface{}, data []byte) interface{} {
	_ = json.Unmarshal(data, &dataStruct)
	return dataStruct
}

func Validations(g *gin.Context, vali interface{}) error {
	if err := validations.Validate.Struct(vali); err != nil {
		return validations.Translate(err.(validator.ValidationErrors))
	}

	return nil
}
