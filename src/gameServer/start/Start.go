package main

import (
	"gameServer"
	"os"
)

func testMain() {
	logDir := "../../logs"
	serverId := "games1001"
	env := "ly"
	gameServer.ServerInstance = &gameServer.GameServer{
		HandlePlayerMap: make(map[int32]gameServer.HandlePlayerFunc),
	}
	gameServer.ServerInstance.StartServer(logDir, serverId, env, initHandler)
}

func main() {
	gameServer.ServerInstance = &gameServer.GameServer{
		HandlePlayerMap: make(map[int32]gameServer.HandlePlayerFunc),
	}

	logDir := ""
	serverId := ""
	env := ""
	if len(os.Args) > 1 {
		logDir = os.Args[1]
		serverId = os.Args[2]
		env = os.Args[3]
	}
	gameServer.ServerInstance.StartServer(logDir, serverId, env, initHandler)
}

func initHandler() {
	gameServer.Inithandler(&gameServer.LoginHandler{})
}
