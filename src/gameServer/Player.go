package gameServer

import (
	"common/utils"
	"server"
	"sync"
	"table"
	"time"
)

// game player client 心跳间隔 100 毫秒
const GameMaxTickTimer = 100

type Player struct {
	Client        *GameClient
	PlayerId      int64
	Name          string
	isStart       bool
	lock          sync.Mutex
	lastTickTimer int64 //上次心跳时间
	PlayerData    table.PlayerDataTableProxy
}

func NewPlayer(playerId int64, client server.NetClientInterface) *Player {
	gameClient := client.(*GameClient)
	return &Player{Client: gameClient, PlayerId: playerId, lastTickTimer: utils.GetNow()}
}

// StartRun 启动玩家逻辑 处理网络包 处理各种事件 100 毫秒心跳一次
func (this *Player) StartRun() {
	if !this.isStart {
		this.isStart = true
	}
	go func() {
		for {
			now := utils.GetNow()
			//处理所有 网络包
			this.Client.TickNet(ServerInstance.ConnectManger)

			//TODO 跨天检查
			if !utils.IsSameDay(now, this.lastTickTimer) {
				//  抛出跨天事件
			}
			//TODO 跨周检查
			this.lastTickTimer = now
			time.Sleep(GameMaxTickTimer * time.Millisecond)
		}
	}()
}
