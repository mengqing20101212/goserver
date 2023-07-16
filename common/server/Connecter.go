package server

import (
	"bytes"
	"fmt"
	"goserver/common/logger"
	"goserver/common/utils"
	"net"
)

func CreateConnect(addr string) *SocketChannel {
	con, err := net.Dial("tcp", addr)
	if err != nil {
		return nil
	}
	tcpConn, ok := con.(*net.TCPConn)
	if !ok {
		logger.Error(fmt.Sprintf("tcpConn, ok := con.(*net.TCPConn) fail addr:%s, con:%s", addr, con))
		return nil
	}
	sc := SocketChannel{
		endPoint:   tcpConn.RemoteAddr(),
		cid:        0,
		con:        tcpConn,
		outputData: make(chan []byte, 64),
		inputMsg:   utils.NewByteBufferByBuf(bytes.NewBuffer(make([]byte, DefaultInputLen))),
	}
	go sc.sendMsg()
	return &sc
}
