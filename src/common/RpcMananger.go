package common

import (
	"common/utils"
	"fmt"
	"github.com/golang/protobuf/proto"
	"protobufMsg"
	"server"
	"sync"
	"time"
)

type RpcMessageInterFace interface {
	GetRpcMessageStructMsgId() int32
}

const MaxMessageSendTimeout = 5 * time.Second

var ServerNodeMap = make(map[string]*ServerNode)
var ServerNodeMapLock = &sync.RWMutex{}

type RpcNode struct {
	ServerNode
	isCreate bool // 是否创建socket
	connet   *server.Connector
	lock     sync.Mutex
}

// SyncSendStructMsg 同步发送结构体消息
// 该函数用于向指定的目标同步发送一个结构体消息。
// 返回值为错误类型，如果发送成功则返回nil，否则返回相应的错误信息。
func (this *RpcNode) SyncSendStructMsg(data RpcMessageInterFace, result *any) error {
	if !this.isCreate {
		err := this.tryConnect()
		if err != nil {
			return err
		}
	}
	dataArr, err := utils.Struct2Bytes(data)
	if err != nil {
		return err
	}
	message := &protobufMsg.CsServer2Server{
		Type: protobufMsg.ServerMsgType_structMsg,
		Cmd:  data.GetRpcMessageStructMsgId(),
		Data: dataArr,
	}
	response, err := this.SyncSendProtoMsg(int32(protobufMsg.CMD_Server2Server), message)
	if err != nil {
		return err
	}
	return utils.Bytes2Struct(response.(*protobufMsg.ScServer2Server).Data, result)
}

// SyncSendProtoMsg 同步发送proto消息
// 该函数用于向指定的目标同步发送一个proto消息。
// 返回值为错误类型，如果发送成功则返回nil，否则返回相应的错误信息。
func (this *RpcNode) SyncSendProtoMsg(cmd int32, data proto.Message) (proto.Message, error) {
	if !this.isCreate {
		err := this.tryConnect()
		if err != nil {
			return nil, err
		}
	}
	flag, pack := this.connet.SendMsgData(cmd, data, MaxMessageSendTimeout)
	if flag {
		responseMsg := protobufMsg.CreateProtoResponseMessage(pack.Cmd)
		if responseMsg == nil {
			return nil, fmt.Errorf("not found response message by cmd:%d", pack.Cmd)
		}
		responseMsg.Reset()
		err := proto.UnmarshalMerge(pack.Body, responseMsg)
		if err != nil {
			return nil, err
		}
		return responseMsg, nil
	} else {
		this.isCreate = false
		return nil, fmt.Errorf("send msg error")
	}
}

func (this *RpcNode) AsyncSendProtoMessage(cmd int32, data proto.Message, callBack server.CallBackFunc) error {
	if !this.isCreate {
		err := this.tryConnect()
		if err != nil {
			return err
		}
	}
	if this.connet.AsyncSendMsgData(cmd, data, MaxMessageSendTimeout, callBack) {
		return nil
	} else {
		return fmt.Errorf("send msg error")
	}
}

func (this *RpcNode) AsyncSendStructMessage(req RpcMessageInterFace, result any, callBack func(bool, any)) error {
	if !this.isCreate {
		err := this.tryConnect()
		if err != nil {
			return err
		}
	}
	dataArr, err := utils.Struct2Bytes(req)
	if err != nil {
		return err
	}
	message := &protobufMsg.CsServer2Server{
		Type: protobufMsg.ServerMsgType_structMsg,
		Cmd:  req.GetRpcMessageStructMsgId(),
		Data: dataArr,
	}
	if this.connet.AsyncSendMsgData(int32(protobufMsg.CMD_Server2Server), message, MaxMessageSendTimeout, func(successFlag bool, responsePacket *server.Package) {
		if successFlag {
			err := utils.Bytes2Struct(responsePacket.Body, &result)
			if err != nil {
				callBack(false, nil)
			} else {
				callBack(true, result)
			}
		} else {
			callBack(false, nil)
		}
	}) {
		return nil
	} else {
		return fmt.Errorf("send msg error")
	}
}

func (this *RpcNode) tryConnect() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.connet = server.CreateConnect(fmt.Sprintln(this.GetIP(), ":", this.ServerPort), server.GeneralCodec)
	err := this.connet.Reconnect()
	this.isCreate = true
	return err
}

func RegisterServerNode(serverNode *ServerNode) {
	ServerNodeMapLock.Lock()
	defer ServerNodeMapLock.Unlock()
	ServerNodeMap[serverNode.ServerId] = serverNode
}
func UnRegisterServerNode(serverId string) {
	ServerNodeMapLock.Lock()
	defer ServerNodeMapLock.Unlock()
	delete(ServerNodeMap, serverId)
}

func getServerNode(serverId string) *ServerNode {
	ServerNodeMapLock.RLock()
	defer ServerNodeMapLock.RUnlock()
	return ServerNodeMap[serverId]
}

func GetServerNodeMapByServerType(serverType ServerType) map[string]*ServerNode {
	ServerNodeMapLock.RLock()
	defer ServerNodeMapLock.RUnlock()
	resultMap := make(map[string]*ServerNode)
	for k, v := range ServerNodeMap {
		if v.ServerType == serverType {
			resultMap[k] = v
		}
	}
	return resultMap
}
