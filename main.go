package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ray1422/dcard-backend-2023/controller"
)

func main() {
	r := gin.Default()
	controller.Register(r.Group("/list"))
	r.Run()
}
