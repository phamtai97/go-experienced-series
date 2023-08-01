package router

import (
	"github.com/gin-gonic/gin"
)

func Post(router *gin.RouterGroup, path string, handler func(c *gin.Context)) {
	router.POST(path, handler)
}
