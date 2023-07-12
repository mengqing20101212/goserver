package server

import (
	"bytes"
	"errors"
	"fmt"
)

type CodeProto interface {
	Decoder(buffer bytes.Buffer) (packageMsg Package, err error)
	Encode(cmd int, traceId int, sendTimer int, sid int, body bytes.Buffer) (packData []byte)
}
type Package struct {
	packageLen uint16
	cmd        int
	sendTimer  int
	traceId    int
	sid        uint16
	bodyLen    uint16
	body       []byte
}

const PackageDefaultHeadLen = 20

func (self Package) Decoder(buffer bytes.Buffer) (packageMsg Package, err error) {

	return Package{}, errors.New("decoder error")
}

func (self Package) Encode(cmd int, traceId int, sendTimer int, sid uint16, body bytes.Buffer) (packData []byte) {
	pack := Package{cmd: cmd, sendTimer: sendTimer, sid: sid, body: body, bodyLen: uint16(body.Len()), traceId: traceId}
	pack.packageLen = self.bodyLen + PackageDefaultHeadLen
	return pack
}

func (p Package) String() string {
	return fmt.Sprintf("{packageLen:%d,cmd:%d,sendTimer:%d,traceId:%d,sid:%d,bodyLen:%d}", p.packageLen, p.cmd, p.sendTimer, p.traceId, p.sid, p.bodyLen)
}
