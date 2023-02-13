package main

import (
	"net"

	"github.com/gin-gonic/gin"
	"github.com/ray1422/dcard-backend-2023/controller"
	"github.com/ray1422/dcard-backend-2023/utils"
)

func main() {
	r := gin.Default()
	lis, err := net.Listen("tcp", utils.Getenv("GRPC_ADDR", "0.0.0.0:3000"))
	if err != nil {
		panic(err)
	}
	controller.Register(r.Group("/list"), lis)
	r.Run()
}
