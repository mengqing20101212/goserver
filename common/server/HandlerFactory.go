package server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"goserver/common/logger"
	"goserver/protobuf"
)

type MsgHandlerInterface interface {
	Execute(msg proto.Message, channel *SocketChannel) (res bool, message proto.Message)
}

type LoginHandler struct {
}

func (self *LoginHandler) Execute(msg proto.Message, channel *SocketChannel) (res bool, message proto.Message) {
	req, ok := msg.(*protobuf.CsLogin)
	if !ok {
		return false, nil
	}
	logger.Info(fmt.Sprintf("req:%s", req))
	return false, &*req
}

func CreateHandler(pack *Package) MsgHandlerInterface {
	switch pack.cmd {
	case 100:
		return &LoginHandler{}

	}
	return nil
}

func CreateProtoMessage(cmd int32) (msg proto.Message) {
	switch cmd {
	case 100:
		return &protobuf.CsLogin{}

	}
	return nil
}
