package router

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/controllers"
	"goPanel/src/panel/middlewares"
)

type Route struct {
	g *gin.Engine
}

func (r *Route) Init(gin *gin.Engine) *gin.Engine {
	r.g = gin

	r.loadGlobalMiddleware()
	r.loadRoute()

	return r.g
}

func (r *Route) loadGlobalMiddleware() {
	r.g.Use(
		new(middlewares.CoreMiddleware).Middleware(),
	)
}

var (
	userController    = controllers.NewUserController()
	machineController = controllers.NewMachineController()
	wsController      = controllers.NewWsController()
)

func (r *Route) loadRoute() {
	r.g.POST("/login", userController.Login)
	r.g.POST("/userAdd", userController.UserAdd)

	r.g.Use(new(middlewares.TokenMiddleware).Middleware())
	routeIndex := r.g.Group("/index")
	{
		routeIndex.GET("/index", controllers.Index)
	}

	// websocket
	routeWs := r.g.Group("/ws")
	{
		routeWs.GET("/ssh/:cols/:rows/:host", wsController.Ssh)
	}

	routeMachine := r.g.Group("/machine")
	{
		routeMachine.GET("/list/:page", machineController.List)
		routeMachine.POST("/add", machineController.Add)
		routeMachine.POST("/edit", machineController.Edit)
		routeMachine.POST("/del", machineController.Del)
	}
}
