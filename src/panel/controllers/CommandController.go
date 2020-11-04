package controllers

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/common"
	"goPanel/src/panel/constants"
)

type CommandController struct {
	BaseController
}

func NewCommandController() *CommandController {
	return &CommandController{}
}

func (c *CommandController) Add(g *gin.Context) {
	common.RetJson(g, constants.SUCCESS, constants.SUCCESS_MSG, "")
}
