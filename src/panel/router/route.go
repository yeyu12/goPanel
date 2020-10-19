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
	routeIndex := r.g.Group("/index")
	{
		routeIndex.GET("/index", controllers.Index)
	}

	// websocket
	routeWs := r.g.Group("/ws")
	{
		routeWs.GET("/ssh/:cols/:rows/:host", controllers.Ssh)
	}

	routeMachine := r.g.Group("/machine")
	{
		routeMachine.GET("/list/:page", controllers.MachineList)
		routeMachine.POST("/add", controllers.MachineAdd)
		routeMachine.POST("/edit", controllers.MachineEdit)
		routeMachine.POST("/del", controllers.MachineDel)
	}
}
