package server

import (
	"bytes"
	"common/utils"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

type Connector struct {
	SocketChannel
	protoCode CodeProto[Package]
	traceId   int32
	sid       uint16
}

func (this *Connector) Send(bs []byte) {
	this.SendMsg(bs)
	this.traceId++
}

func (this *Connector) SendMsgData(cmd int32, msg proto.Message) (flag bool, responsePack *Package) {
	if !this.IsConnect() {
		this.reconnect()
		if !this.IsConnect() {
			return false, nil
		}
	}
	responseData, err := proto.Marshal(msg)
	if err != nil {
		log.Error(fmt.Sprintf("up pack msg error:%s, cmd:%d, msg:%s, endPoint:%s", err, cmd, msg.String(), this.endPoint.String()))
		this.Close(" Connector close by marshal data error ")
		return false, nil
	}
	pack := CreatePackage(int32(cmd), this.traceId, uint32(time.Now().Unix()), this.sid, responseData)
	bs := this.protoCode.Encode(pack)
	this.Send(bs)
	readLen, err := this.con.Read(this.inputMsg.GetBytes())
	if err != nil {
		log.Error(fmt.Sprintf("read romote data error:%s, cmd:%d, reqPack:%s，endPoint:%s", err, cmd, pack, this.endPoint.String()))
		this.Close(" Connector close by read remote data error ")
		return false, nil
	}
	if readLen > 0 {
		responsePack, unpackFlag := this.protoCode.Decoder(&this.inputMsg)
		if unpackFlag {
			if responsePack.sid > 0 && this.sid == 0 {
				this.sid = responsePack.sid
				log.Info(fmt.Sprintf("set new sid:%d, endPoint:%s", this.sid, this.endPoint.String()))
			}
			return true, responsePack
		}
	}
	return false, nil
}

func (this *Connector) reconnect() {
	if this.con != nil {
		this.con.Close()
	}
	con, err := net.Dial("tcp", this.socketIp)
	if err != nil {
		log.Error(fmt.Sprintf("reconnect connect endPoint:%s fail, err:%s", this.socketIp, err))
		return
	}
	this.con = con.(*net.TCPConn)
	this.cid = 1
	this.endPoint = this.con.RemoteAddr()
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
