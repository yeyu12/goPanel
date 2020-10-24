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

func (r *Route) loadRoute() {
	r.g.POST("/login", controllers.NewUserController().Login)
	r.g.POST("/userAdd", controllers.NewUserController().UserAdd)

	// websocket
	routeWs := r.g.Group("/ws")
	{
		routeWs.GET("/ssh", controllers.NewWsController().Ssh)
	}

	r.g.Use(new(middlewares.TokenMiddleware).Middleware())
	routeIndex := r.g.Group("/index")
	{
		routeIndex.GET("/index", controllers.Index)
	}

	routeMachine := r.g.Group("/machine")
	{
		routeMachine.GET("/list", controllers.NewMachineController().List)
		routeMachine.POST("/save", controllers.NewMachineController().Save)
		routeMachine.POST("/del", controllers.NewMachineController().Del)
	}
}
