package controllers

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/common"
	"goPanel/src/panel/constants"
	core "goPanel/src/panel/core/database"
	"goPanel/src/panel/models"
	"goPanel/src/panel/services"
	"goPanel/src/panel/validations"
	"io/ioutil"
	"time"
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

	if err := c.Validations(userVail); err != nil {
		common.RetJson(g, constants.ERROR_FAIL, err.Error(), "")
		return
	}
	var postMap map[string]string
	c.JsonPost(&postMap, inputData)

	postMap["passwd"] = common.StringUtils(common.StringUtils(postMap["passwd"]).SHA1()).MD5()
	data := c.userService.UsernameAndPasswdByData(core.Db, postMap)

	if data.Id == 0 {
		common.RetJson(g,
			constants.USERNAME_OR_PASSWD_ERROR_FAIL,
			constants.USERNAME_OR_PASSWD_ERROR_FAIL_MSG,
			"")

		return
	}

	token, err := common.GetToken()
	if err != nil {
		common.RetJson(g, constants.ERROR_FAIL, constants.ERROR_FAIL_MSG, "")
		return
	}

	// 更新数据
	data.Token = token
	data.TokenExpirationTime = time.Now()

	affected, err := c.userService.UpdateUser(core.Db, data)
	if affected == 0 || err != nil {
		common.RetJson(g, constants.ERROR_FAIL, constants.ERROR_FAIL_MSG, "")
		return
	}

	data.Passwd = ""

	common.RetJson(g,
		constants.SUCCESS,
		constants.SUCCESS_MSG, data)

	return
}

func (c *UserController) UserAdd(g *gin.Context) {
	inputData, _ := ioutil.ReadAll(g.Request.Body)
	var userVail validations.UserAdd
	c.JsonPost(&userVail, inputData)

	if err := c.Validations(userVail); err != nil {
		common.RetJson(g, constants.ERROR_FAIL, err.Error(), "")
		return
	}
	var userAddData models.UserModel
	c.JsonPost(&userAddData, inputData)

	if oldUserData := c.userService.UsernameData(core.Db, userAddData.Username); oldUserData.Id != 0 {
		common.RetJson(g, constants.USERNAME_ALREADY_EXISTS, constants.USERNAME_ALREADY_EXISTS_MSG, oldUserData)
		return
	}

	token, err := common.GetToken()
	if err != nil {
		common.RetJson(g, constants.ERROR_FAIL, constants.ERROR_FAIL_MSG, "")
		return
	}

	// 构建数据
	userAddData.Token = token
	userAddData.Passwd = common.StringUtils(common.StringUtils(userAddData.Passwd).SHA1()).MD5()
	userAddData.CreateTime = time.Now()
	userAddData.TokenExpirationTime = time.Now()

	userId, err := c.userService.UserAdd(core.Db, userAddData)
	if err != nil || userId == 0 {
		common.RetJson(g,
			constants.ERROR_FAIL,
			constants.ERROR_FAIL_MSG,
			"")

		return
	}

	userData := c.userService.UsernameAndPasswdByData(core.Db, map[string]string{
		"username": userAddData.Username,
		"passwd":   userAddData.Passwd,
	})

	common.RetJson(g,
		constants.SUCCESS,
		constants.SUCCESS_MSG, userData)

	return
}
