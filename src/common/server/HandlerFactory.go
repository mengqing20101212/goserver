package server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"protobufMsg"
)

type MsgHandlerInterface interface {
	Execute(msg proto.Message, channel *SocketChannel) (res bool, message proto.Message)
	HandLerName() string
}

var handlerMap = make(map[int32]MsgHandlerInterface)

type EmptyHandler struct {
}

func (resf *EmptyHandler) Execute(msg proto.Message, channel *SocketChannel) (res bool, message proto.Message) {
	log.Error("not run EmptyHandler error")
	return false, nil
}
func (self *EmptyHandler) HandLerName() string {
	return "EmptyHandler"
}

type LoginHandler struct{ EmptyHandler }

func (self *LoginHandler) HandLerName() string {
	return "LoginHandler"
}

func CreateHandler(pack *Package) MsgHandlerInterface {
	handler := handlerMap[pack.cmd]
	if handler != nil {
		return handler
	}
	log.Error(fmt.Sprintf(" not found msg Handler cmd:%d, pack:%s", pack.cmd, pack))
	return nil
}

func InitHandler(cmd int32, handler MsgHandlerInterface) {
	handlerMap[cmd] = handler
	log.Info(fmt.Sprintf(" add new cmdId:%d, handler:%s", cmd, handler.HandLerName()))
}

func DefaultInitHandler() {

	InitHandler(100, &LoginHandler{})

}

func CreateProtoRequestMessage(cmd int32) (msg proto.Message) {
	switch cmd {

	case 100:
		return &protobufMsg.CsLogin{}

	default:
		log.Error(fmt.Sprintf(" not found cmdId:%d", cmd))
		return nil

	}
	log.Error(fmt.Sprintf(" not found cmdId:%d", cmd))
	return nil
}
