package server

import (
	"bytes"
	"fmt"
	"goserver/common/logger"
	"goserver/common/utils"
)

type CodeProto interface {
	Decoder(buffer bytes.Buffer) (packageMsg *Package, success bool)
	Encode() (packData []byte)
}
type Package struct {
	packageLen uint16
	cmd        int32
	sendTimer  uint32
	traceId    int32
	sid        uint16
	bodyLen    uint16
	body       []byte
}

const PackageDefaultHeadLen = 2 + 4 + 4 + 4 + 2 + 2 // packageLen（2） + cmd（4）+ sendTimer（4）+traceId（4）+ sid（2） + bodyLen（2）

func (self *Package) Decoder(buffer bytes.Buffer) (packageMsg *Package, success bool) {
	buf := utils.NewByteBufferByBuf(&buffer)
	buf.Mark()
	packLen := buf.ReadUint16()
	if len(buf.GetBytes()) < int(packLen) {
		buf.RestMark()
		return nil, false
	}
	cmd := buf.ReadInt32()
	sendTimer := buf.ReadUint32()
	traceId := buf.ReadInt32()
	sid := buf.ReadUint16()
	bodyLen := buf.ReadUint16()
	body := buf.ReadBytes(int(bodyLen))
	msg := CreatePackage(cmd, traceId, sendTimer, sid, body)
	return msg, true
}

func (self *Package) Encode() (packData []byte) {
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
func CreatePackage(cmd int32, traceId int32, sendTimer uint32, sid uint16, body []byte) (packData *Package) {
	pack := Package{cmd: cmd, sendTimer: sendTimer, sid: sid, body: body, bodyLen: uint16(len(body)), traceId: traceId}
	pack.packageLen = PackageDefaultHeadLen + pack.bodyLen
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("new package:%s", pack.String()))
	}
	return &pack
}
func (self *Package) String() string {
	return fmt.Sprintf("{packageLen:%d,cmd:%d,sendTimer:%d,traceId:%d,sid:%d,bodyLen:%d}", self.packageLen, self.cmd, self.sendTimer, self.traceId, self.sid, self.bodyLen)
}
