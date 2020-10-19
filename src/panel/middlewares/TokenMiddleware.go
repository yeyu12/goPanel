package middlewares

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/common"
)

type TokenMiddleware struct {
}

func (core *TokenMiddleware) Middleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		//token := g.Request.Header.Get("Account-Token")

		common.RetJson(g, 3000, "请先登录！", "")

		// 处理请求
		g.Next()
	}
}
