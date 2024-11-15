package server

import (
	"common/utils"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type CodeProto[T Package] interface {
	Decoder(buf *utils.ByteBuffer) (packageMsg *T, success bool)
	Encode(msg *T) (packData []byte)
}

type Package struct {
	packageLen uint16
	Cmd        int32
	SendTimer  uint32
	TraceId    int32
	Sid        uint16
	seq        uint32
	bodyLen    uint16
	body       []byte
}

type PackageMessage struct {
	*Package
	proto.Message
}

const PackageDefaultHeadLen = 2 + 4 + 4 + 4 + 2 + 4 + 2 // packageLen（2） + cmd（4）+ sendTimer（4）+traceId（4）+ sid（2）+ seq(4) + bodyLen（2）

type PackageFactory struct {
}

func (self *PackageFactory) Decoder(buf *utils.ByteBuffer) (packageMsg *Package, success bool) {
	buf.Mark()
	packLen, err := buf.ReadUint16()
	if err != nil {
		buf.RestMark()
		return nil, false
	}
	byteLen := buf.Len()
	if byteLen+2 < int(packLen) {
		buf.RestMark()
		return nil, false
	}
	cmd := buf.ReadInt32()
	sendTimer := buf.ReadUint32()
	traceId := buf.ReadInt32()
	sid, _ := buf.ReadUint16()
	bodyLen, _ := buf.ReadUint16()
	_ = buf.ReadInt32()
	body := buf.ReadBytes(int(bodyLen))
	msg := CreatePackage(cmd, traceId, sendTimer, sid, body)
	return msg, true
}

func (self *PackageFactory) Encode(msg *Package) (packData []byte) {
	packBuf := utils.NewByteBuffer()
	packBuf.WriteUint16(msg.packageLen)
	packBuf.WriteInt32(msg.Cmd)
	packBuf.WriteUint32(msg.SendTimer)
	packBuf.WriteInt32(msg.TraceId)
	packBuf.WriteUint16(msg.Sid)
	packBuf.WriteUint16(msg.bodyLen)
	packBuf.WriteInt32(int32(msg.seq))
	packBuf.WriteBytes(msg.body)
	return packBuf.GetBytes()
}
func CreatePackage(cmd int32, traceId int32, sendTimer uint32, sid uint16, body []byte) (packData *Package) {
	pack := Package{Cmd: cmd, SendTimer: sendTimer, Sid: sid, body: body, bodyLen: uint16(len(body)), TraceId: traceId}
	pack.packageLen = PackageDefaultHeadLen + pack.bodyLen
	if log.IsDebug() {
		log.Debug(fmt.Sprintf("new package:%s", pack.String()))
	}
	return &pack
}
func (self *Package) String() string {
	return fmt.Sprintf("{packageLen:%d,cmd:%d,sendTimer:%d,traceId:%d,sid:%d,bodyLen:%d,body:%s}", self.packageLen, self.Cmd, self.SendTimer, self.TraceId, self.Sid, self.bodyLen, self.body)
}
