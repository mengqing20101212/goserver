package main

import "sync"

type Player struct {
	client   *GameClient
	PlayerId int64
	Name     string
	isStart  bool
	lock     sync.Mutex
}

func NewPlayer(playerId int64, client *GameClient) *Player {
	return &Player{client: client, PlayerId: playerId}
}

func (this *Player) StartRun() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if !this.isStart {
		this.isStart = true
	}

}
