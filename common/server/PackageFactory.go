package server

import (
	"fmt"
	"goserver/common/logger"
	"goserver/common/utils"
)

type CodeProto interface {
	Decoder(buf *utils.ByteBuffer) (packageMsg *Package, success bool)
	Encode() (packData []byte)
}
type Package struct {
	packageLen uint16
	cmd        int32
	sendTimer  uint32
	traceId    int32
	sid        uint16
	seq        uint32
	bodyLen    uint16
	body       []byte
}

const PackageDefaultHeadLen = 2 + 4 + 4 + 4 + 2 + 4 + 2 // packageLen（2） + cmd（4）+ sendTimer（4）+traceId（4）+ sid（2）+ seq(4) + bodyLen（2）

type PackageFactory struct {
	Package
}

func (self *PackageFactory) Decoder(buf *utils.ByteBuffer) (packageMsg *Package, success bool) {
	buf.Mark()
	packLen, err := buf.ReadUint16()
	if err != nil {
		buf.RestMark()
		return nil, false
	}
	byteLen := buf.GetBuffer().Len()
	if byteLen < int(packLen) {
		buf.RestMark()
		return nil, false
	}
	cmd := buf.ReadInt32()
	sendTimer := buf.ReadUint32()
	traceId := buf.ReadInt32()
	sid, _ := buf.ReadUint16()
	bodyLen, _ := buf.ReadUint16()
	body := buf.ReadBytes(int(bodyLen))
	msg := CreatePackage(cmd, traceId, sendTimer, sid, body)
	return &msg.Package, true
}

func (self *PackageFactory) Encode() (packData []byte) {
	packBuf := utils.NewByteBuffer()
	packBuf.WriteUint16(self.packageLen)
	packBuf.WriteInt32(self.cmd)
	packBuf.WriteUint32(self.sendTimer)
	packBuf.WriteInt32(self.traceId)
	packBuf.WriteUint16(self.sid)
	packBuf.WriteUint16(self.bodyLen)
	packBuf.WriteBytes(self.body)
	return packBuf.GetBytes()
}
func CreatePackage(cmd int32, traceId int32, sendTimer uint32, sid uint16, body []byte) (packData *PackageFactory) {
	pack := PackageFactory{Package{cmd: cmd, sendTimer: sendTimer, sid: sid, body: body, bodyLen: uint16(len(body)), traceId: traceId}}
	//pack := Package{cmd: cmd, sendTimer: sendTimer, sid: sid, body: body, bodyLen: uint16(len(body)), traceId: traceId}
	pack.packageLen = PackageDefaultHeadLen + pack.bodyLen
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("new package:%s", pack.String()))
	}
	return &pack
}
func (self *Package) String() string {
	return fmt.Sprintf("{packageLen:%d,cmd:%d,sendTimer:%d,traceId:%d,sid:%d,bodyLen:%d,body:%s}", self.packageLen, self.cmd, self.sendTimer, self.traceId, self.sid, self.bodyLen, self.body)
}
