package gameServer

import (
	"common"
	"common/utils"
	"fmt"
	"github.com/golang/protobuf/proto"
	"logger"
	"protobufMsg"
	"server"
)

var gameLogger *logger.Logger

type HandlePlayerFunc func(msg proto.Message, channel *GameClient, player *Player) (res bool, responseMessage proto.Message)

type GameServer struct {
	server.Server
	handlePlayerMap map[int32]HandlePlayerFunc
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
	} else { //玩家已经登录
		player := PlayerManagerInstance.GetPlayer(this.playerId)
		handle := ServerInstance.handlePlayerMap[cmd]
		if handle != nil {
			returnFlag, response := handle(data.PackageMessage.Message, this, player)
			responseData, err := proto.Marshal(response)
			if err != nil {
				this.CloseNet(fmt.Sprintf(" parse receivePackageMessage response req:%s, response:%s, pack:%s, error:%s", data.PackageMessage.Message.String(), response, data.PackageMessage.Package.String(), err), mgr)
				return true
			}
			resPack := server.CreatePackage(data.PackageMessage.Cmd, data.PackageMessage.TraceId, data.PackageMessage.SendTimer, data.PackageMessage.Sid, responseData)
			this.SocketChannel.SendMsg(server.GeneralCodec.Encode(resPack))
			return returnFlag
		} else {
			this.CloseNet(fmt.Sprintf(" unknown cmd: %d, sid: %d", cmd, this.GetCid()), mgr)
		}
	}
	return false
}

var ServerInstance GameServer

func (this *GameServer) StartServer(serverId, env string) {
	begin := utils.GetNow()
	ServerInstance.Server = server.NewServer(common.Context.Config.ServerPort)
	gameLogger = common.InitContext("../logs", serverId, env, common.Game, &ServerInstance)
	gameLogger.Info("GameServer InitContext success")
	ServerInstance.Start()
	server.CreateServerStatus(&ServerInstance.Server, common.Game.String(), serverId, common.ServerRunModule.String())
	gameLogger.Info(fmt.Sprintf("GameServer StartServer success  useCost:%d", (utils.GetNow()-begin)/1000))
}

func Inithandler(handler server.HandlerInterface) {
	handler.Initializer()
}

func AddGameServerHandler(cmd int32, handle HandlePlayerFunc) {
	ServerInstance.handlePlayerMap[cmd] = handle

}

func (this *GameServer) StopServer() {
	//TODO  保存数据库数据
	common.CloseContext()
	gameLogger.Info("GameServer StopServer success")
}

func GetConnectManger() *server.ConnectManger {
	return ServerInstance.ConnectManger

}
