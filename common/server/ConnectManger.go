package server

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"goserver/common/logger"
	"goserver/common/protos"
)

type ConnectManger struct {
	connMap map[int32]*SocketChannel
}

func (self *ConnectManger) AddConn(channel *SocketChannel, server *Server) {
	if self.connMap[channel.cid] != nil {
		logger.Error(fmt.Sprintf("[ConnectManger] repeact add SocketChannel:%s", channel))
		return
	}
	self.connMap[channel.cid] = channel
	logger.Info(fmt.Sprintf("[ConnectManger] AddConn:%s", channel))
	go lookupReadData(channel, server)
}

func lookupReadData(channel *SocketChannel, server *Server) {
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
			cslogin := &protos.CsLogin{}
			proto.UnmarshalMerge(pack.body, cslogin)

		}

	}
}

func (self *ConnectManger) DelConn(channel *SocketChannel) {
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
