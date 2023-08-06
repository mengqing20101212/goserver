package handler

import (
	"github.com/gin-gonic/gin"
	"log"
)

func TestGet(c *gin.Context) {
	log.Default().Println(c.RemoteIP())
	c.JSON(200, gin.H{
		"remoter": c.RemoteIP(),
	})
}
