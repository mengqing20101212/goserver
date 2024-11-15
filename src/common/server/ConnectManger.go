package server

import (
	"common/utils"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"net"
	"protobufMsg"
	"sync"
	"time"
)

type ConnectManger struct {
	connMap map[uint16]NetClientInterface
	lock    sync.Mutex
}

const maxReceivePackageMessageLen = 1024
const tickSleepTimer = 100 * time.Millisecond //心跳间隔 毫秒

type MessageDataType int

const (
	receivePackageMessage MessageDataType = iota //PackageMessage消息
	sendMessage                                  //protobuf 消息
)

type NetClientInterface interface {
	IsRunning() bool
	AddReceiveMsg(packet *Package, msg proto.Message)
	AddSendMsg(cmd int32, msg proto.Message)
	CloseNet(closeMsg string, mgr *ConnectManger)
	SendMsg(data []byte)
	TickNet()
	HandleReceivePackageMessage(data *OptionData, mgr *ConnectManger) bool
	GetCid() uint16
	GetSocketChannel() *SocketChannel
}

// 待处理消息
type OptionData struct {
	optType        MessageDataType
	PackageMessage *PackageMessage
	Message        proto.Message //发送的消息
	sendCmdId      int32         //发送的消息号
}

type NetClient struct {
	*SocketChannel
	lock    sync.Mutex
	msgList utils.List[*OptionData] //需要处理的包
	start   bool
}

func NewNetClient(channel *SocketChannel) *NetClient {
	return &NetClient{
		SocketChannel: channel,
		msgList:       utils.NewList[*OptionData](),
		start:         false,
	}
}

func (receiver *NetClient) GetSocketChannel() *SocketChannel {
	return receiver.SocketChannel

}

func (this *NetClient) GetCid() uint16 {
	return this.cid
}

func (this *NetClient) GetConn() *net.TCPConn {
	return this.con
}

func (this *NetClient) SendMsg(data []byte) {
	this.SocketChannel.SendMsg(data)
}

func (this *NetClient) startRun(mgr *ConnectManger) {
	this.start = true
	for {
		//检查当前socket
		if !this.IsRunning() {
			this.CloseNet(fmt.Sprintf("NetClient [%s], colse ", this.endPoint.String()), mgr)
			return
		}

		//处理子类消息
		this.TickNet()

		// 处理消息
		this.msgList.ForEachAndClear(func(data *OptionData) {
			this.lock.Lock()
			defer this.lock.Unlock()
			switch data.optType {

			case receivePackageMessage: //收到远端传来的消息
				if this.HandleReceivePackageMessage(data, mgr) {
					return
				}

			case sendMessage: //其他的携程写入到 发送队列的消息
				responseData, err := proto.Marshal(data.Message)
				if err != nil {
					this.CloseNet(fmt.Sprintf(" parse sendMessage response req:%s, response:%s,  error:%s", data.PackageMessage.Message.String(), data.Message.String(), err), mgr)
					return
				}
				resPack := CreatePackage(data.sendCmdId, 0, uint32(time.Now().Unix()), this.cid, responseData)
				this.SocketChannel.SendMsg(GeneralCodec.Encode(resPack))
			}
		})
		//休息一下
		time.Sleep(time.Duration(tickSleepTimer))
	}
}

func (this *NetClient) HandleReceivePackageMessage(data *OptionData, mgr *ConnectManger) bool {
	handler := CreateHandler(data.PackageMessage.Package)
	returnFlag, response := handler(data.PackageMessage.Message, this)
	if !returnFlag {
		logger.Info(fmt.Sprintf(" package no result req:%s, pack:%s", data.Message.String(), data.PackageMessage.Package.String()))
		return true
	}
	responseData, err := proto.Marshal(response)
	if err != nil {
		this.CloseNet(fmt.Sprintf(" parse receivePackageMessage response req:%s, response:%s, pack:%s, error:%s", data.PackageMessage.Message.String(), response, data.PackageMessage.Package.String(), err), mgr)
		return true
	}
	resPack := CreatePackage(data.PackageMessage.Cmd, data.PackageMessage.TraceId, data.PackageMessage.SendTimer, data.PackageMessage.Sid, responseData)
	this.SocketChannel.SendMsg(GeneralCodec.Encode(resPack))
	return false
}

func (this *NetClient) CloseNet(closeMsg string, mgr *ConnectManger) {
	this.msgList.Clear()
	this.con.Close()
	mgr.DelConn(this)
	log.Info(closeMsg)
}
func (this *NetClient) IsRunning() bool {
	return this.IsConnect()
}
func (this *NetClient) TickNet() {

}

func (this *NetClient) AddReceiveMsg(packet *Package, msg proto.Message) {
	packetMessage := &PackageMessage{
		packet,
		msg,
	}
	data := &OptionData{
		optType:        receivePackageMessage,
		PackageMessage: packetMessage,
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.msgList.Size() >= maxReceivePackageMessageLen {
		log.Error(fmt.Sprintf("too many packet curLen:%d", this.msgList.Size()))
		return
	}
	this.msgList.Add(data)
}

func (this *NetClient) AddSendMsg(cmd int32, msg proto.Message) {
	data := &OptionData{
		optType:   sendMessage,
		Message:   msg,
		sendCmdId: cmd,
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.msgList.Size() >= maxReceivePackageMessageLen {
		log.Error(fmt.Sprintf("too many packet curLen:%d", this.msgList.Size()))
		return
	}
	this.msgList.Add(data)
}

func (this *NetClient) isStart() bool {
	return this.start
}

type ServerNetClient struct {
	NetClient
}

func (this *ServerNetClient) HandleReceivePackageMessage(data *OptionData, mgr *ConnectManger) bool {
	return true
}

// 从socket 读取数据 并分发到指定的client 处理
func loopReadData(channel NetClientInterface, server *Server, mgr *ConnectManger) {
	for {
		bs := make([]byte, 256)
		n, err := channel.GetSocketChannel().inputMsg.GetBuffer().Read(bs)
		if err != nil {
			channel.CloseNet(fmt.Sprintf("[lookupReadData] read data error:%s, channel:%s", err, channel), mgr)
			return
		}
		if n <= 0 {
			continue
		}
		channel.GetSocketChannel().inputMsg.GetBuffer().Write(bs[:n])
		for {
			pack, res := GeneralCodec.Decoder(&channel.GetSocketChannel().inputMsg)
			if res {
				if !server.filterChain.Filter(pack, channel.GetSocketChannel()) {
					channel.CloseNet(fmt.Sprintf("[lookupReadData]  close channel by Filter pack:%s, channel:%s", pack, channel), mgr)
					return
				}
			} else {
				break
			}
			reqMessage := protobufMsg.CreateProtoRequestMessage(pack.Cmd)
			if reqMessage == nil {
				channel.CloseNet(fmt.Sprintf("Message not found pack:%s", pack), mgr)
				return
			}
			reqMessage.Reset()
			err = proto.UnmarshalMerge(pack.body, reqMessage)
			if err != nil {
				logger.Error()
				channel.CloseNet(fmt.Sprintf("UnmarshalMerge pack:%s  error:%s", pack, err), mgr)
				return
			}
			channel.AddReceiveMsg(pack, reqMessage)

		}

	}
}

func (this *ConnectManger) AddConn(channel NetClientInterface, server *Server) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.connMap[channel.GetCid()] != nil {
		logger.Error(fmt.Sprintf("[ConnectManger] repeact add SocketChannel:%s", channel))
		return
	}
	this.connMap[channel.GetCid()] = channel
	logger.Info(fmt.Sprintf("[ConnectManger] AddConn:%s", channel))
	go loopReadData(channel, server, this)
}

func (this *ConnectManger) DelConn(channel *NetClient) {
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
		connMap: make(map[uint16]NetClientInterface, maxConLen),
	}
	return &manger
}
