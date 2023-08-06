package main

import (
	"github.com/gin-gonic/gin"
	"goserver/LoginServer/handler"
)

func Test2(c *gin.Context) {

}

func RegisterGet(router *gin.Engine) {
	router.GET("test/1/", handler.TestGet)
	router.GET("test/2/", Test2)
}
