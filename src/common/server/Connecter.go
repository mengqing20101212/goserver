package server

import (
	"bytes"
	"common/utils"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"protobuf"
	"time"
)

type Connector struct {
	SocketChannel
	protoCode CodeProto[Package]
	traceId   int32
	sid       uint16
}

func (self *Connector) Send(bs []byte) {
	self.SendMsg(bs)
	self.traceId++
}

func (self *Connector) SendMsgData(cmd protobuf.CMD, msg proto.Message) (flag bool, responsePack *Package) {
	if !self.IsConnect() {
		self.reconnect()
		if !self.IsConnect() {
			return false, nil
		}
	}
	responseData, err := proto.Marshal(msg)
	if err != nil {
		log.Error(fmt.Sprintf("up pack msg error:%s, cmd:%d, msg:%s, endPoint:%s", err, cmd, msg.String(), self.endPoint.String()))
		self.Close(" Connector close by marshal data error ")
		return false, nil
	}
	pack := CreatePackage(int32(cmd), self.traceId, uint32(time.Now().Unix()), self.sid, responseData)
	bs := self.protoCode.Encode(pack)
	self.Send(bs)
	readLen, err := self.con.Read(self.inputMsg.GetBytes())
	if err != nil {
		log.Error(fmt.Sprintf("read romote data error:%s, cmd:%d, reqPack:%s，endPoint:%s", err, cmd, pack, self.endPoint.String()))
		self.Close(" Connector close by read remote data error ")
		return false, nil
	}
	if readLen > 0 {
		responsePack, unpackFlag := self.protoCode.Decoder(&self.inputMsg)
		if unpackFlag {
			if responsePack.sid > 0 && self.sid == 0 {
				self.sid = responsePack.sid
				log.Info(fmt.Sprintf("set new sid:%d, endPoint:%s", self.sid, self.endPoint.String()))
			}
			return true, responsePack
		}
	}
	return false, nil
}

func (self *Connector) reconnect() {
	if self.con != nil {
		self.con.Close()
	}
	con, err := net.Dial("tcp", self.socketIp)
	if err != nil {
		log.Error(fmt.Sprintf("reconnect connect endPoint:%s fail, err:%s", self.socketIp, err))
		return
	}
	self.con = con.(*net.TCPConn)
	self.cid = 1
	self.endPoint = self.con.RemoteAddr()
}

func CreateConnect(addr string, protoCode CodeProto[Package]) *Connector {
	sc := Connector{
		SocketChannel: SocketChannel{
			cid:      -1,
			socketIp: addr,
			inputMsg: utils.NewByteBufferByBuf(bytes.NewBuffer(make([]byte, DefaultInputLen))),
		},
		protoCode: protoCode,
		traceId:   1,
	}
	sc.reconnect()
	return &sc
}
