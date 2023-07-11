package server

import (
	"bytes"
	"goserver/common/logger"
)

type CodeProto interface {
	Decoder(buffer bytes.Buffer) (packageMsg Package, err error)
	Encode(body bytes.Buffer) (msg Package)
}
type Package struct {
	packageLen uint16
	cmd        int
	sendTimer  int
	traceId    int
	sid        uint16
	bodyLen    uint16
	body       bytes.Buffer
}

type FilterInterface interface {
	DoFilter(msg *SocketChannel) (success bool)
}
type FilterChain struct {
	next *FilterInterface
}

func (self FilterChain) Filter(channel *SocketChannel) (success bool) {
	if self.next != nil {
		b := (*self.next).DoFilter(channel)
		if !b {
			logger.Error("")
			return b
		}
	}
	return true
}
func (self FilterChain) AddFilter(filter *FilterInterface) {
	self.next = filter
}

/*
type FilterChain struct {
	next *Filter
}*/
