package server

import (
	"bytes"
	"fmt"
	"gameProject/common/utils"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

type Connector struct {
	SocketChannel
	protoCode CodeProto[Package]
	traceId   int32
}

func (this *Connector) Send(bs []byte) {
	this.SendMsg(bs)
	this.traceId++
}

type CallBackFunc func(successFlag bool, responsePacket *Package)

func (this *Connector) AsyncSendMsgData(cmd int32, msg proto.Message, timeOut time.Duration, callBack CallBackFunc) bool {
	pack, _, _, done := this.sendData(cmd, msg)
	if done {
		return false
	}
	go this.handRemoteData(timeOut, cmd, pack, callBack)
	return true
}

func (this *Connector) handRemoteData(timeOut time.Duration, cmd int32, pack any, callBack CallBackFunc) {
	for {
		time.Sleep(100 * time.Millisecond)
		timeOut -= 100 * time.Millisecond
		if timeOut <= 0 {
			log.Error(fmt.Sprintf("AsyncSendMsgData send msg timeout, cmd:%d, reqPack:%s，endPoint:%s", cmd, pack, this.endPoint.String()))
			this.Close("AsyncSendMsgData Connector close by send msg timeout ")
			callBack(false, nil)
		}

		readLen, err := this.con.Read(this.inputMsg.GetBytes())
		if err != nil {
			log.Error(fmt.Sprintf("read romote data error:%s, cmd:%d, reqPack:%s，endPoint:%s", err, cmd, pack, this.endPoint.String()))
			this.Close(" Connector close by read remote data error ")
			callBack(false, nil)
		}
		if readLen > 0 {
			responsePack, unpackFlag := this.protoCode.Decoder(&this.inputMsg)
			if unpackFlag {
				if responsePack.Sid > 0 && this.cid == 0 {
					this.cid = responsePack.Sid
					log.Info(fmt.Sprintf("set new sid:%d, endPoint:%s", this.cid, this.endPoint.String()))
				}
				callBack(true, responsePack)
			}
		}
	}
}

// SendMsgData sends a proto.Message to the server with the given command and processes the response.
// Returns a boolean flag indicating success, and the response package if successful.
func (this *Connector) SendMsgData(cmd int32, msg proto.Message, timeOut time.Duration) (flag bool, responsePack *Package) {
	pack, b, p, done := this.sendData(cmd, msg)
	if done {
		return b, p
	}
	for {
		time.Sleep(100 * time.Millisecond)
		timeOut -= 100 * time.Millisecond
		if timeOut <= 0 {
			log.Error(fmt.Sprintf("send msg timeout, cmd:%d, reqPack:%s，endPoint:%s", cmd, pack, this.endPoint.String()))
			this.Close(" Connector close by send msg timeout ")
			return false, nil
		}
		readLen, err := this.con.Read(this.inputMsg.GetBytes())
		if err != nil {
			log.Error(fmt.Sprintf("read romote data error:%s, cmd:%d, reqPack:%s，endPoint:%s", err, cmd, pack, this.endPoint.String()))
			this.Close(" Connector close by read remote data error ")
			return false, nil
		}
		if readLen > 0 {
			responsePack, unpackFlag := this.protoCode.Decoder(&this.inputMsg)
			if unpackFlag {
				if responsePack.Sid > 0 && this.cid == 0 {
					this.cid = responsePack.Sid
					log.Info(fmt.Sprintf("set new sid:%d, endPoint:%s", this.cid, this.endPoint.String()))
				}
				return true, responsePack
			}
		}
		return false, nil
	}
}

func (this *Connector) sendData(cmd int32, msg proto.Message) (*Package, bool, *Package, bool) {
	if !this.IsConnect() {
		this.Reconnect()
		if !this.IsConnect() {
			return nil, false, nil, true
		}
	}
	responseData, err := proto.Marshal(msg)
	if err != nil {
		log.Error(fmt.Sprintf("up pack msg error:%s, cmd:%d, msg:%s, endPoint:%s", err, cmd, msg.String(), this.endPoint.String()))
		this.Close(" Connector close by marshal data error ")
		return nil, false, nil, true
	}
	pack := CreatePackage(int32(cmd), this.traceId, uint32(time.Now().Unix()), this.cid, responseData)
	bs := this.protoCode.Encode(pack)
	this.Send(bs)
	return pack, false, nil, false
}

func (this *Connector) Reconnect() error {
	if this.con != nil {
		this.con.Close()
	}
	con, err := net.Dial("tcp", this.socketIp)
	if err != nil {
		log.Error(fmt.Sprintf("reconnect connect endPoint:%s fail, err:%s", this.socketIp, err))
		return err
	}
	this.con = con.(*net.TCPConn)
	// 读写缓冲区都为 10K 不够在扩大
	err = this.con.SetReadBuffer(1024 * 10)
	err = this.con.SetWriteBuffer(1024 * 10)
	if err != nil {
		return err
	}
	this.cid = 0
	this.endPoint = this.con.RemoteAddr()
	return nil
}

func CreateConnect(addr string, protoCode CodeProto[Package]) *Connector {
	sc := Connector{
		SocketChannel: SocketChannel{
			cid:      0,
			socketIp: addr,
			inputMsg: utils.NewByteBufferByBuf(bytes.NewBuffer(make([]byte, DefaultInputLen))),
		},
		protoCode: protoCode,
		traceId:   1,
	}
	sc.Reconnect()
	return &sc
}
