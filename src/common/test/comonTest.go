package main

import (
	"common"
	"common/utils"
	"encoding/base64"
	"logger"
	"time"
)

func main() {
	//testLogger()
	testConfig()
	//testNacos()

}

func testConfig() {
	common.InitContext("../../logs", "game1001", common.Game)
}

func testNacos() {
	utils.InitNacos("../../logs", "game1001", nil)
}

func base64Str(str string) string {
	bs := []byte(str)
	return base64.StdEncoding.EncodeToString(bs)
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
