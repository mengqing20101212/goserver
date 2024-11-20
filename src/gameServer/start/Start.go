package main

import "gameServer"

func main() {
	gameServer := gameServer.GameServer{}
	gameServer.StartServer("game1001", "ly")
}
