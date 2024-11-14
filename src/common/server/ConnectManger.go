package server

import (
	"common/utils"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"protobufMsg"
	"sync"
	"time"
)

type ConnectManger struct {
	connMap map[uint16]*NetClient
	lock    sync.Mutex
}

const maxReceivePackageMessageLen = 1024
const tickSleepTimer = 100 * time.Millisecond //心跳间隔 毫秒

type MessageDataType int

const (
	receivePackageMessage MessageDataType = iota //PackageMessage消息
	sendMessage                                  //protobuf 消息
)

type NetStateCheck interface {
	IsRunning() bool
}

// 待处理消息
type OptionData struct {
	optType        MessageDataType
	packageMessage *PackageMessage
	message        proto.Message //发送的消息
	sendCmdId      int32         //发送的消息号
}

type NetClient struct {
	SocketChannel
	lock    sync.Mutex
	msgList utils.List[*OptionData] //需要处理的包
	start   bool
}

func (this *NetClient) startRun(mgr *ConnectManger) {
	this.start = true
	for {
		//检查当前socket
		if !this.IsRunning() {
			log.Info(fmt.Sprintf("NetClient [%s], colse ", this.endPoint.String()))
			this.CloseNet(mgr)
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
				handler := CreateHandler(data.packageMessage.Package)
				returnFlag, response := handler(&data.packageMessage.Message, &this.SocketChannel)
				if !returnFlag {
					logger.Info(fmt.Sprintf(" package no result req:%s, pack:%s", data.message.String(), data.packageMessage.Package.String()))
					return
				}
				responseData, err := proto.Marshal(response)
				if err != nil {
					logger.Error(fmt.Sprintf(" parse receivePackageMessage response req:%s, response:%s, pack:%s, error:%s", data.packageMessage.Message.String(), response, data.packageMessage.Package.String(), err))
					this.CloseNet(mgr)
					return
				}
				resPack := CreatePackage(data.packageMessage.cmd, data.packageMessage.traceId, data.packageMessage.sendTimer, data.packageMessage.sid, responseData)
				this.SocketChannel.SendMsg(GeneralCodec.Encode(resPack))

			case sendMessage: //其他的携程写入到 发送队列的消息
				responseData, err := proto.Marshal(data.message)
				if err != nil {
					logger.Error(fmt.Sprintf(" parse sendMessage response req:%s, response:%s,  error:%s", data.packageMessage.Message.String(), data.message.String(), err))
					this.CloseNet(mgr)
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

func (this *NetClient) CloseNet(mgr *ConnectManger) {
	this.msgList.Clear()
	this.con.Close()
	mgr.DelConn(this)
}
func (this *NetClient) IsRunning() bool {
	return this.IsConnect()
}
func (this *NetClient) TickNet() {

}

func NewNetClient(channel SocketChannel) *NetClient {
	netClient := &NetClient{
		SocketChannel: channel,
		lock:          sync.Mutex{},
		msgList:       utils.NewList[*OptionData](),
	}
	return netClient
}
func (this *NetClient) AddReceiveMsg(packet *Package, msg proto.Message) {
	packetMessage := &PackageMessage{
		packet,
		msg,
	}
	data := &OptionData{
		optType:        receivePackageMessage,
		packageMessage: packetMessage,
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
		message:   msg,
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

// 从socket 读取数据 并分发到指定的client 处理
func loopReadData(channel *NetClient, server *Server, mgr *ConnectManger) {
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
			pack, res := GeneralCodec.Decoder(&channel.inputMsg)
			if res {
				if !server.filterChain.Filter(pack, &channel.SocketChannel) {
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
			channel.AddReceiveMsg(pack, reqMessage)
			if !channel.isStart() {
				channel.startRun(mgr)
			}
		}

	}
}

func (this *ConnectManger) AddConn(channel *NetClient, server *Server) {
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
		connMap: make(map[uint16]*NetClient, maxConLen),
	}
	return &manger
}
