package main

import (
	"goserver/common/logger"
	"time"
)

var log logger.Logger

func main() {
	//testLogger()
}

func testLogger() {
	for i := 0; i < 10; i++ {
		go func() {
			log := logger.InitNull()
			log.Debug("debug")
			log.Info("info")
			log.Warn("warn")
			log.Error("error")
		}()
	}
	time.Sleep(10 * time.Second)
}
