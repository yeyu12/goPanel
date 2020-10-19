package controllers

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/common"
	"goPanel/src/panel/constants"
	core "goPanel/src/panel/core/database"
	"goPanel/src/panel/services"
	"goPanel/src/panel/validations"
	"io/ioutil"
)

type UserController struct {
	BaseController
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: new(services.UserService),
	}
}

func (c *UserController) Login(g *gin.Context) {
	inputData, _ := ioutil.ReadAll(g.Request.Body)
	var userVail validations.Login
	c.JsonPost(&userVail, inputData)

	if err := c.Validations(g, userVail); err != nil {
		common.RetJson(g, constants.ConstantsCode["ERROR_FAIL"], err.Error(), "")
		return
	}
	var postMap map[string]string
	c.JsonPost(&postMap, inputData)

	data := c.userService.UsernameAndPasswdByData(core.Db, postMap)

	if data.Id == 0 {
		common.RetJson(g,
			constants.ConstantsCode["USERNAME_OR_PASSWD_ERROR_FAIL"],
			constants.ConstantsMsg["USERNAME_OR_PASSWD_ERROR_FAIL"],
			"")
	}

	common.RetJson(g,
		constants.ConstantsCode["SUCCESS"],
		constants.ConstantsMsg["SUCCESS"], data)

	return
}
