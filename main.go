package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ray1422/dcard-backend-2023/controllers"
)

func main() {
	r := gin.Default()
	controllers.Register(r.Group("/list"))
	r.Run()
}
