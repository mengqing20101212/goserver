package server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"goserver/common/logger"
	"goserver/protobuf"
)

type MsgHandlerInterface interface {
	Execute(msg proto.Message) (res bool, message proto.Message)
}

type BaseHandler struct {
	pack      *Package
	channel   *SocketChannel
	protoCode CodeProto[Package]
}

func (self *BaseHandler) SendMsg(msg proto.Message) {
	resBytes, err := proto.Marshal(msg)
	if err != nil {
		self.channel.Close(fmt.Sprintf(" unpack msg error:%s  pack:%s", msg, self.pack))
	}
	resPack := CreatePackage(self.pack.cmd, self.pack.traceId, self.pack.sendTimer, uint16(self.channel.cid), resBytes)
	self.channel.SendMsg(self.protoCode.Encode(resPack))
}

type LoginHandler struct {
	BaseHandler
}

func (self *LoginHandler) Execute(msg proto.Message) (res bool, message proto.Message) {
	req, ok := msg.(*protobuf.CsLogin)
	if !ok {
		return false, nil
	}
	self.SendMsg(&*req)
	logger.Info(fmt.Sprintf("req:%s", req))
	return false, &*req
}

func CreateHandler(pack *Package, channel *SocketChannel, protoCode CodeProto[Package]) MsgHandlerInterface {
	hanler := LoginHandler{
		BaseHandler{
			pack:      pack,
			channel:   channel,
			protoCode: protoCode,
		},
	}
	return &hanler
}
