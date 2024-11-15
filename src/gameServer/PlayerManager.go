package main

import (
	"fmt"
	"sync"
)

type PlayerManager struct {
	olinePlayerMap map[int64]*Player
	lock           sync.Mutex
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		olinePlayerMap: make(map[int64]*Player),
	}
}

func (self *PlayerManager) AddPlayer(player *Player) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.olinePlayerMap[player.PlayerId] = player
}

func (self *PlayerManager) DelPlayer(playerId int64) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.olinePlayerMap, playerId)
}

func (self *PlayerManager) GetPlayer(playerId int64) *Player {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.olinePlayerMap[playerId]
}

func (self *PlayerManager) ForeachPlayer(foreachFunc func(player *Player)) {
	for _, player := range self.olinePlayerMap {
		foreachFunc(player)
	}
}

type Test interface {
	TestLog()
}
type BaseT struct {
	i1 int
}

func (this *BaseT) TestLog() {
	fmt.Println("BaseT log")
}

type DerivedT struct {
	BaseT
	i2 int
}

func (this *DerivedT) TestLog() {
	fmt.Println("DerivedT log")
}

func main() {
	m := make(map[string]Test)
	base := BaseT{}
	derived := DerivedT{}
	m["base"] = &base
	m["derived"] = &derived
	for _, v := range m {
		v.TestLog()
	}
}
