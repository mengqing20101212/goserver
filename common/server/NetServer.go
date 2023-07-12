package server

import "fmt"

type Server struct {
	port        int
	proto       string
	codecsProto *CodeProto
	filterChain *FilterChain
}

type EndPoint struct {
	ip   string
	port uint16
}

func (e EndPoint) String() string {
	return fmt.Sprintf("{ip:%s,port:%d}", e.ip, e.port)
}

func NewEndPoint(ip string, port uint16) *EndPoint {
	return &EndPoint{ip: ip, port: port}
}

type SocketChannel struct {
	msg      Package
	endPoint EndPoint
}

func (e SocketChannel) String() string {
	return fmt.Sprintf("SocketChannel{msg:%s,endPoint:%d}", e.msg.String(), e.endPoint.String())
}

func NewSocketChannel(msg Package, endPoint EndPoint) *SocketChannel {
	return &SocketChannel{msg: msg, endPoint: endPoint}
}
