package gameServer

import (
	"github.com/golang/protobuf/proto"
	"protobufMsg"
	"server"
)

type LoginHandler struct {
}

func (receiver *LoginHandler) Initializer() {
	server.InitHandler(int32(protobufMsg.CMD_cmd_login), login)
	server.InitHandler(int32(protobufMsg.CMD_cmd_login3), login3)
	server.InitHandler(int32(protobufMsg.CMD_cmd_login2), login2)
	server.InitHandler(int32(protobufMsg.CMD_cmd_login4), login4)
}

func (receiver *LoginHandler) HandleName() string {
	return "loginHandler"
}

func login(msg *proto.Message, channel *server.SocketChannel) (res bool, responseMessage proto.Message) {
	return false, nil
}

func login2(msg *proto.Message, channel *server.SocketChannel) (res bool, responseMessage proto.Message) {
	return false, nil
}
func login3(msg *proto.Message, channel *server.SocketChannel) (res bool, responseMessage proto.Message) {
	return false, nil
}
func login4(msg *proto.Message, channel *server.SocketChannel) (res bool, responseMessage proto.Message) {
	return false, nil
}
