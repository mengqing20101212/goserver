package main

import (
	"encoding/base64"
	"gameProject/common"
	"gameProject/common/utils"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"logger"
	"os"
	"time"
)

func main() {
	//testLogger()
	if len(os.Args) > 1 {
		testConfig(os.Args[1])
	} else {
		testConfig("game1001")
	}
	//testNacos()

}

func testConfig(serverId string) {
	common.InitContext("../../logs", serverId, "ly", common.Game, nil)
	time.Sleep(1023 * time.Second)
}

func testNacos() {
	utils.InitNacos("../../logs", "game1001", "ly", nil)
	utils.RegisterNewServerCallBack(common.Game.String(), func(isAdd bool, instance model.Instance) {
	})
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
