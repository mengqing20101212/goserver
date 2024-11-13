package server

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"protobufMsg"
	"sync"
)

type ConnectManger struct {
	connMap map[uint16]*SocketChannel
	lock    sync.Mutex
}

const maxReceivePackageMessageLen = 1024
const maxSendPackageMessageLen = 10

type MessageDataType int

const (
	packageMessage MessageDataType = iota //PackageMessage消息
	message                               //protobuf 消息
)

type OptionData struct {
	optType        MessageDataType
	packageMessage *PackageMessage
	message        proto.Message
}

type NetClient struct {
	SocketChannel
	lock            sync.Mutex
	receiveMsgQueue chan *PackageMessage //收到远端的包
	sendMsgQueue    chan proto.Message   //待发送的远端的包
}

func (this *NetClient) tick() {

}
func (this *NetClient) TickNet() {

}

func NewNetClient(channel SocketChannel) *NetClient {
	netClient := &NetClient{
		SocketChannel:   channel,
		lock:            sync.Mutex{},
		receiveMsgQueue: make(chan *PackageMessage, maxReceivePackageMessageLen),
		sendMsgQueue:    make(chan proto.Message, maxSendPackageMessageLen),
	}
	return netClient
}
func (this *NetClient) AddReceiveMsg(packet *Package, msg proto.Message) {
	packetMessage := &PackageMessage{
		packet,
		msg,
	}
	this.receiveMsgQueue <- packetMessage
}
func (this *NetClient) AddSendMsg(packet proto.Message) {
	this.sendMsgQueue <- packet
}

// 从socket 读取数据 并分发到指定的client 处理
func loopReadData(channel *SocketChannel, server *Server, mgr *ConnectManger) {
	for {
		bs := make([]byte, 256)
		n, err := channel.con.Read(bs)
		if err != nil {
			channel.Close(fmt.Sprintf("[lookupReadData] read data error:%s, channel:%s", err, channel))
			server.ConnectManger.DelConn(channel)
			return
		}
		if n <= 0 {
			continue
		}
		channel.inputMsg.GetBuffer().Write(bs[:n])
		for {
			pack, res := server.codecsProto.Decoder(&channel.inputMsg)
			if res {
				if !server.filterChain.Filter(pack, channel) {
					channel.Close(fmt.Sprintf("[lookupReadData]  close channel by Filter pack:%s, channel:%s", pack, channel))
					server.ConnectManger.DelConn(channel)
					return
				}
			} else {
				break
			}
			reqMessage := protobufMsg.CreateProtoRequestMessage(pack.cmd)
			if reqMessage == nil {
				logger.Error(fmt.Sprintf("Message not found pack:%s", pack))
				mgr.DelConn(channel)
				return
			}
			reqMessage.Reset()
			err = proto.UnmarshalMerge(pack.body, reqMessage)
			if err != nil {
				logger.Error(fmt.Sprintf("UnmarshalMerge pack:%s  error:%s", pack, err))
				mgr.DelConn(channel)
				return
			}

			handler := CreateHandler(pack)
			returnFlag, response := handler(&reqMessage, channel)
			if !returnFlag {
				logger.Info(fmt.Sprintf(" package no result req:%s, pack:%s", reqMessage.String(), pack.String()))
				continue
			}
			responseData, err := proto.Marshal(response)
			if err != nil {
				logger.Error(fmt.Sprintf(" parse response req:%s, response:%s, pack:%s, error:%s", reqMessage.String(), response, pack.String(), err))
				mgr.DelConn(channel)
				return
			}
			resPack := CreatePackage(pack.cmd, pack.traceId, pack.sendTimer, pack.sid, responseData)
			channel.SendMsg(server.codecsProto.Encode(resPack))
		}

	}
}

func (this *ConnectManger) AddConn(channel *SocketChannel, server *Server) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.connMap[channel.cid] != nil {
		logger.Error(fmt.Sprintf("[ConnectManger] repeact add SocketChannel:%s", channel))
		return
	}
	this.connMap[channel.cid] = channel
	logger.Info(fmt.Sprintf("[ConnectManger] AddConn:%s", channel))
	go loopReadData(channel, server, this)
}

func (this *ConnectManger) DelConn(channel *SocketChannel) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.connMap, channel.cid)
	logger.Info(fmt.Sprintf("[ConnectManger] delete SocketChannel:%s", channel))
}

func (this *ConnectManger) SendMsgToConn(cid uint16, sendData []byte) error {
	channel := this.connMap[cid]
	if channel == nil {
		logger.Error(fmt.Sprintf("[SendMsgToConn] not found channel:%d", cid))
		return errors.New("SendMsgToConn not found channel cid:" + string(cid))
	}
	channel.SendMsg(sendData)
	return nil
}

func NewConnectManger(maxConLen int) (mgr *ConnectManger) {
	manger := ConnectManger{
		connMap: make(map[uint16]*SocketChannel, maxConLen),
	}
	return &manger
}
