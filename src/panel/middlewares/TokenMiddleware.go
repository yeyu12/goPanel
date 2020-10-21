package middlewares

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/common"
	"goPanel/src/panel/services"
)

type TokenMiddleware struct {
	userService *services.UserService
}

func (m *TokenMiddleware) Middleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		m.userService = new(services.UserService)
		token := g.Request.Header.Get("Account-Token")

		state, msg, code := m.userService.IsUserLogin(token)
		if !state {
			common.RetJson(g, code, msg, "")
			return
		}

		// 处理请求
		g.Next()
	}
}
