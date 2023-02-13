package controller

import (
	"net"

	"github.com/gin-gonic/gin"
)

// Register routers under this module
func Register(r *gin.RouterGroup, listener net.Listener) {
	r.GET("/:id/:version", listHandler)
	r.GET("/:id", headHandler)
	go func() {
		err := listenGRPC(listener)
		if err != nil {
			panic(err)
		}
	}()

}
