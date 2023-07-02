package main

import (
	"fmt"
	"goserver/common/logger"
)

var log logger.Logger

func main() {
	dir := "../logs"
	filename := "test.log"
	log.Init(&dir, &filename)
	logConfig := logger.LoggerConfig{}
	log.InitByConfig(logConfig)
	fmt.Println("12")
	log.Debug("dada")
}
