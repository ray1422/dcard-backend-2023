package controller

import "github.com/gin-gonic/gin"

// Register routers under this module
func Register(r *gin.RouterGroup) {
	r.GET("/:id/:version", listHandler)
	r.GET("/:id", headHandler)
}
