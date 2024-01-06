package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := InitWebookServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello，启动成功了！")
	})
	server.Run(":8080")
}
