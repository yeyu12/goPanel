package main

import (
	"github.com/gin-gonic/gin"
	"goPanel/src/panel/router"
)

func main() {
	g := gin.Default()
	g = (new(router.Route)).Init(g)
	_ = g.Run(":10010") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
