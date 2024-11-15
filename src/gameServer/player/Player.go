package gameServer

import (
	"gameServer"
	"server"
	"sync"
)

type Player struct {
	Client   *gameServer.GameClient
	PlayerId int64
	Name     string
	isStart  bool
	lock     sync.Mutex
}

func NewPlayer(playerId int64, client server.NetClientInterface) *Player {
	gameClient := client.(*gameServer.GameClient)
	return &Player{Client: gameClient, PlayerId: playerId}
}

func (this *Player) StartRun() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if !this.isStart {
		this.isStart = true
	}
	go func() {

	}()
}
