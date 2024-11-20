package main

import (
	"gameServer"
	gameServer2 "gameServer/login"
)

func main() {

	gameServer := gameServer.GameServer{}
	gameServer.StartServer("game1001", "ly")
	initHandler()

}

func initHandler() {
	loginHandler := gameServer2.LoginHandler{}
	loginHandler.Initializer()

}
