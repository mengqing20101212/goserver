package main

import "github.com/gin-gonic/gin"

func main() {
	route := gin.Default()

	route.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	RegisterGet(route)
	route.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
