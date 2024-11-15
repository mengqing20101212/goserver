package gameServer

import (
	"common"
	"fmt"
	gameServer "gameServer/login"
	gameServer2 "gameServer/player"
	"github.com/golang/protobuf/proto"
	"logger"
	"server"
)

var gameLogger *logger.Logger

type GameServer struct {
	server.Server
}

func (this *GameServer) CreateNewClient(channel *server.SocketChannel) server.NetClientInterface {
	gameClient := &GameClient{
		NetClient: *server.NewNetClient(channel),
	}
	return gameClient
}

type GameClient struct {
	server.NetClient
	playerId int64
}

func (this *GameClient) TickNet() {

}
func (this *GameClient) HandleReceivePackageMessage(data *server.OptionData, mgr *server.ConnectManger) bool {
	handler := server.CreateHandler(data.PackageMessage.Package)
	returnFlag, response := handler(data.PackageMessage.Message, this)
	if !returnFlag {
		gameLogger.Info(fmt.Sprintf(" package no result req:%s, pack:%s", data.Message.String(), data.PackageMessage.Package.String()))
		return true
	}
	responseData, err := proto.Marshal(response)
	if err != nil {
		this.CloseNet(fmt.Sprintf(" parse receivePackageMessage response req:%s, response:%s, pack:%s, error:%s", data.PackageMessage.Message.String(), response, data.PackageMessage.Package.String(), err), mgr)
		return true
	}
	resPack := server.CreatePackage(data.PackageMessage.Cmd, data.PackageMessage.TraceId, data.PackageMessage.SendTimer, data.PackageMessage.Sid, responseData)
	this.SocketChannel.SendMsg(server.GeneralCodec.Encode(resPack))
	return false
}

var GameServerInstance GameServer
var PlayerManger = gameServer2.NewPlayerManager()

func (this *GameServer) StartServer(serverId, env string) {

	common.InitContext("../logs", serverId, env, common.Game)
	gameLogger = logger.Init(common.Context.Config.LogDir, "game")
	this.InitHandler()
	gameLogger.Info("GameServer InitContext success")
	GameServerInstance.Server = server.NewServer(common.Context.Config.ServerPort)
	GameServerInstance.Start()
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
