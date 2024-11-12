package main

import (
	"common"
	gameServer "gameServer/login"
	"logger"
	"server"
)

type GameServer struct {
}

var gameLogger *logger.Logger

func (this *GameServer) StartServer(serverId, env string) {

	common.InitContext("../logs", serverId, env, common.Game)
	gameLogger = logger.Init(common.Context.Config.LogDir, "game")
	this.InitHandler()
	gameLogger.Info("GameServer StartServer success")

}

func (this *GameServer) InitHandler() {
	initHandler(&gameServer.LoginHandler{})
}

func initHandler(handler server.HandlerInterface) {
	handler.Initializer()
}

func (this *GameServer) StopServer() {
	//TODO  保存数据库数据
	gameLogger.Info("GameServer StopServer success")
}
