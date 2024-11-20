package gameServer

import (
	gameServer "gameServer/player"
	"github.com/golang/protobuf/proto"
	"protobufMsg"
	"server"
)

type LoginHandler struct {
}

func (receiver *LoginHandler) Initializer() {
	server.InitHandler(int32(protobufMsg.CMD_Login), login)
	server.InitHandler(int32(protobufMsg.CMD_Login3), login3)
	server.InitHandler(int32(protobufMsg.CMD_Login2), login2)
}

func (receiver *LoginHandler) HandleName() string {
	return "loginHandler"
}

func login(msg proto.Message, channel server.NetClientInterface) (res bool, responseMessage proto.Message) {
	req := msg.(*protobufMsg.CsLogin)

	newPlayer := gameServer.NewPlayer(int64(req.Scores[0]), channel)
	//gameServer2.PlayerManger.AddPlayer(newPlayer)
	newPlayer.StartRun()
	return false, &protobufMsg.ScLogin{
		Name:   req.Name,
		Male:   req.Male,
		Scores: req.Scores,
	}
}

func login2(msg proto.Message, channel server.NetClientInterface) (res bool, responseMessage proto.Message) {
	return false, nil
}
func login3(msg proto.Message, channel server.NetClientInterface) (res bool, responseMessage proto.Message) {
	return false, nil
}
func login4(msg *proto.Message, channel server.NetClientInterface) (res bool, responseMessage proto.Message) {
	return false, nil
}
