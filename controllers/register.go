package controllers

import "github.com/gin-gonic/gin"

// Register routers under this module
func Register(r *gin.RouterGroup) {
	r.GET("/:key/:user_id", HeadHandler)
}
