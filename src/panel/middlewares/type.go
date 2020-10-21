package middlewares

import "github.com/gin-gonic/gin"

type MiddlewareInterface interface {
	Middleware() gin.HandlerFunc
}
