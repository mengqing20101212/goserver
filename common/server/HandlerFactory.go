package server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"goserver/common/logger"
	"goserver/protobuf"
)

type MsgHandlerInterface interface {
	Execute(msg interface{}) (res bool)
}

type BaseHandler struct {
	pack      *Package
	channel   *SocketChannel
	protoCode *PackageFactory
}

func (self *BaseHandler) SendMsg(msg *proto.Message) {
	resBytes, err := proto.Marshal(*msg)
	if err != nil {
		self.channel.Close(fmt.Sprintf(" unpack msg error:%s  pack:%s", msg, self.pack))
	}
	resPack := CreatePackage(self.pack.cmd, self.pack.traceId, self.pack.sendTimer, uint16(self.channel.cid), resBytes)
	self.channel.SendMsg(self.protoCode.Encode(resPack))
}

type LoginHandler struct {
	BaseHandler
}

func (self *LoginHandler) Execute(msg interface{}) (res bool) {
	req, ok := msg.(*protobuf.CsLogin)
	if !ok {
		return false
	}
	logger.Info(fmt.Sprintf("req:%s", req))
	return false
}

type HandlerFactory struct {
	handlerMap map[int32]MsgHandlerInterface
}

func (self *HandlerFactory) GetHandler(cmd int32) MsgHandlerInterface {
	return self.handlerMap[cmd]
}
func (self *HandlerFactory) RegisterHandler(cmd int32, handlerInterface MsgHandlerInterface) {
	self.handlerMap[cmd] = handlerInterface
}

var handlerFactory HandlerFactory

func GetHandlerFactoryInstance() *HandlerFactory {
	return &handlerFactory
}
