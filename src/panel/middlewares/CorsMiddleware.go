package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CoreMiddleware struct {
}

func (core *CoreMiddleware) Middleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		method := g.Request.Method

		g.Header("Access-Control-Allow-Origin", "*")
		g.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Account-Token")
		g.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		g.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		g.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			g.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		g.Next()
	}
}
