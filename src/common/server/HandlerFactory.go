package server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"protobufMsg"
)

type HandleFunc func(msg proto.Message, channel NetClientInterface) (res bool, responseMessage proto.Message)

// 初始化 指定模块 并注册所有的 msg handler 处理器
type HandlerInterface interface {
	Initializer()       //初始化模块
	HandleName() string //该模块名称
}

var handlerMap = make(map[int32]HandleFunc)

func CreateHandler(pack *Package) HandleFunc {
	handler := handlerMap[pack.Cmd]
	if handler != nil {
		return handler
	}
	log.Error(fmt.Sprintf(" not found msg Handler cmd:%d, pack:%s", pack.Cmd, pack))
	return nil
}

func InitHandler(cmd int32, handler HandleFunc) {
	handlerMap[cmd] = handler
	log.Info(fmt.Sprintf(" add new cmdId:%d, handler:%s", cmd, protobufMsg.CMD_name[cmd]))
}

func DefaultInitHandler() {

	//InitHandler(100, &LoginHandler{})

}
