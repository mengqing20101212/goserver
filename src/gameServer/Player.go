package gameServer

import (
	"fmt"
	"gameProject/common"
	"gameProject/common/utils"
	"server"
	"sync"
	"table"
	"time"
)

// game player client 心跳间隔 100 毫秒
const GameMaxTickTimer = 100

type Player struct {
	Client        *GameClient
	PlayerId      int64                       //
	Name          string                      // 玩家名字
	isStart       bool                        // 是否启动
	lock          sync.Mutex                  // 锁
	lastTickTimer int64                       //上次心跳时间
	PlayerData    *table.PlayerDataTableProxy // 玩家数据
	eventHandle   EventProcess
	moduleMap     map[ModuleNameEnum]IModule
}
type ModuleNameEnum string

const (
	ModuleNameHero ModuleNameEnum = "HeroModule"
	ModuleNameItem ModuleNameEnum = "ItemModule"
)

type IModule interface {
	InitModule(player *Player)
	GetModuleName() ModuleNameEnum
}

type BaseModule struct {
	Player *Player
}

func (this *BaseModule) InitBaseModule(player *Player) {
	this.Player = player
}

func NewPlayer(playerId int64, client server.NetClientInterface) *Player {
	gameClient := client.(*GameClient)
	newPlayer := &Player{Client: gameClient, PlayerId: playerId, lastTickTimer: utils.GetNow(),
		PlayerData:  table.NewPlayerDataTable(true),
		eventHandle: EventProcess{},
	}
	newPlayer.initAllModules()
	return newPlayer
}

func (this *Player) GetHeroModule() *HeroModule {
	module := this.moduleMap[ModuleNameHero]
	if module != nil {
		return module.(*HeroModule)
	}
	gameLogger.Error(fmt.Sprintf("player %d get module %s failed", this.PlayerId, string(ModuleNameHero)))
	return nil
}

func (this *Player) PostEvent(eventType EventType, params ...any) {
	this.eventHandle.HandleEvent(eventType, params...)
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

			this.PostEvent(EventType_Login, this)
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

func (this *Player) initAllModules() {
	this.moduleMap[ModuleNameHero] = &HeroModule{}

	for enum, module := range this.moduleMap {
		module.InitModule(this)
		if common.IsTest() {
			gameLogger.Debug(fmt.Sprintf("player %d init module %s success", this.PlayerId, string(enum)))
		}
	}
}
