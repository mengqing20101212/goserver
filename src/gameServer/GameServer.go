package gameServer

import (
	"common"
	"common/utils"
	"fmt"
	gameServer "gameServer/login"
	gameServer2 "gameServer/player"
	"logger"
	"protobufMsg"
	"server"
)

var gameLogger *logger.Logger

type GameServer struct {
	server.Server
}

func (this *GameServer) CreateNewClient(channel *server.SocketChannel) server.NetClientInterface {
	gameClient := &GameClient{
		NetClient:     *server.NewNetClient(channel),
		lastTickTimer: utils.GetNow(),
	}
	return gameClient
}

type GameClient struct {
	server.NetClient
	playerId      int64
	lastTickTimer int64
}

func (this *GameClient) HandleReceivePackageMessage(data *server.OptionData, mgr *server.ConnectManger) bool {
	cmd := data.PackageMessage.Cmd
	if this.playerId == 0 {
		if cmd != int32(protobufMsg.CMD_Login) { // 该玩家未登录 非法包
			this.CloseNet(fmt.Sprintln(" 玩家未登录, 不接受其他的包 cmd: ", cmd, ", sid:", this.GetCid()), mgr)
			return true
		}
		return true
	} else {
		this.NetClient.HandleReceivePackageMessage(data, mgr)
	}
	return false
}

var ServerInstance GameServer
var PlayerManger = gameServer2.NewPlayerManager()

func (this *GameServer) StartServer(serverId, env string) {

	common.InitContext("../logs", serverId, env, common.Game)
	gameLogger = logger.Init(common.Context.Config.LogDir, "game")
	this.InitHandler()
	gameLogger.Info("GameServer InitContext success")
	ServerInstance.Server = server.NewServer(common.Context.Config.ServerPort)
	ServerInstance.Start()
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
	this.Server.ConnectManger.CloseAllClient()
	this.Server.Stop()
	gameLogger.Info("GameServer StopServer success")
}

func GetConnectManger() *server.ConnectManger {
	return ServerInstance.ConnectManger

}
