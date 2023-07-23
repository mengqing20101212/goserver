package server

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"goserver/common/logger"
	"sync"
)

type ConnectManger struct {
	connMap map[int32]*SocketChannel
	lock    sync.Mutex
}

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
			reqMessage := CreateProtoRequestMessage(pack.cmd)
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
			returnFlag, response := handler.Execute(reqMessage, channel)
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
			fmt.Println(response)
		}

	}
}

func (self *ConnectManger) AddConn(channel *SocketChannel, server *Server) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.connMap[channel.cid] != nil {
		logger.Error(fmt.Sprintf("[ConnectManger] repeact add SocketChannel:%s", channel))
		return
	}
	self.connMap[channel.cid] = channel
	logger.Info(fmt.Sprintf("[ConnectManger] AddConn:%s", channel))
	go loopReadData(channel, server, self)
}

func (self *ConnectManger) DelConn(channel *SocketChannel) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.connMap, channel.cid)
	logger.Info(fmt.Sprintf("[ConnectManger] delete SocketChannel:%s", channel))
}

func (self *ConnectManger) SendMsgToConn(cid int32, sendData []byte) error {
	channel := self.connMap[cid]
	if channel != nil {
		logger.Error(fmt.Sprintf("[SendMsgToConn] not found channel:%d", cid))
		return errors.New("SendMsgToConn not found channel cid:" + string(cid))
	}
	channel.SendMsg(sendData)
	return nil
}

func NewConnectManger(maxConLen int) (mgr *ConnectManger) {
	manger := ConnectManger{
		connMap: make(map[int32]*SocketChannel, maxConLen),
	}
	return &manger
}
