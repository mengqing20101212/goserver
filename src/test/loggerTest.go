package main

import (
	"goserver/common/logger"
	"time"
)

var log logger.Logger

func main() {
	dir := "/logs"
	filename := "test.log"
	logger.Init(dir, filename)
	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < 1000; i++ {
				logger.Debug("logger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debug")
				logger.Info("logger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debug")
				logger.Error("logger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debug")
				logger.Warn("logger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debuglogger.Debug")
			}
		}()
		time.Sleep(1 * time.Second)
	}
	time.Sleep(10 * time.Second)

}
