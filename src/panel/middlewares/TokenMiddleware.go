package middlewares

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/common"
	"goPanel/src/panel/constants"
)

type TokenMiddleware struct {
}

func (core *TokenMiddleware) Middleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		//token := g.Request.Header.Get("Account-Token")

		common.RetJson(g, constants.PLEASE_LOG_IN_FIRST, constants.PLEASE_LOG_IN_FIRST_MSG, "")

		// 处理请求
		g.Next()
	}
}
