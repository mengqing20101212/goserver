package gameServer

import (
	"github.com/golang/protobuf/proto"
	"protobufMsg"
	"server"
)

type LoginHandler struct {
}

func (receiver *LoginHandler) Initializer() {
	server.InitHandler(int32(protobufMsg.CMD_Login), login)
	AddGameServerHandler(int32(protobufMsg.CMD_Login2), login2)
}

func login2(msg proto.Message, channel *GameClient, player *Player) (res bool, responseMessage proto.Message) {
	return true, nil
}

func (receiver *LoginHandler) HandleName() string {
	return "loginHandler"
}

func login(msg proto.Message, channel server.NetClientInterface) (res bool, responseMessage proto.Message) {
	req := msg.(*protobufMsg.CsLogin)

	newPlayer := NewPlayer(int64(req.Scores[0]), channel)
	newPlayer.PostEvent(EventType_CreateRole, newPlayer)
	PlayerManagerInstance.AddPlayer(newPlayer)
	newPlayer.StartRun()
	return false, &protobufMsg.ScLogin{
		Name:   req.Name,
		Male:   req.Male,
		Scores: req.Scores,
	}
}
