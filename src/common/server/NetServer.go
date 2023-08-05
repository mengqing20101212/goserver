package server

import (
	"bytes"
	"fmt"
	"goserver/common/logger"
	"goserver/common/utils"
	"net"
	"os"
)

const DefaultInputLen = 1024 * 5
const DefaultMaxConnectLen = 1024

type Server struct {
	port          int
	proto         string
	codecsProto   CodeProto[Package]
	filterChain   *FilterChain
	listener      *net.TCPListener
	ConnectManger *ConnectManger
	connectNum    int32
}

func NewServer(port int) (server *Server) {
	server = &Server{
		port:          port,
		proto:         "tcp",
		filterChain:   &FilterChain{},
		ConnectManger: NewConnectManger(DefaultMaxConnectLen),
		connectNum:    1,
	}
	return server
}

func (self *Server) Start() {
	logger.Info("[Server] start begin")
	self.filterChain.AddFilter(&Filter{})
	self.filterChain.AddFilter(&IpFilter{})
	self.codecsProto = &PackageFactory{}
	ipaddr := net.TCPAddr{Port: self.port}
	DefaultInitHandler()
	lis, err := net.ListenTCP("tcp", &ipaddr)
	if err != nil {
		logger.Error(fmt.Sprintf("[Server] start ListenTCP error:%s", err))
		os.Exit(10)
	}
	self.listener = lis
	logger.Info(fmt.Sprintf("[Server] start end listener port:%d", self.port))
	for {
		con, err := self.listener.Accept()
		if err != nil {
			logger.Error(fmt.Sprintf("[Server] Accept error: %s", err))
			continue
		}
		go self.OnAccept(con, self.connectNum)
		self.connectNum++
	}
}

func (self *Server) OnAccept(con net.Conn, cid int32) {
	tcpConn, ok := con.(*net.TCPConn)
	if !ok {
		// 处理类型转换失败的情况
		logger.Error(fmt.Sprintf("OnAccept conn tcpConn, ok := con.(*net.TCPConn) error"))
		return
	}
	addr := tcpConn.RemoteAddr()
	sc := SocketChannel{
		endPoint: addr,
		cid:      cid,
		con:      tcpConn,
		inputMsg: utils.NewByteBufferByBuf(bytes.NewBuffer(make([]byte, DefaultInputLen))),
	}
	sc.inputMsg.GetBuffer().Reset()
	self.ConnectManger.AddConn(&sc, self)
}

type SocketChannel struct {
	endPoint net.Addr
	socketIp string
	cid      int32
	con      *net.TCPConn
	inputMsg utils.ByteBuffer
}

func (e *SocketChannel) String() string {
	return fmt.Sprintf("SocketChannel{endPoint:%s}", e.endPoint.String())
}

func (self *SocketChannel) SendMsg(data []byte) {
	if self.IsConnect() {
		_, err := self.con.Write(data)
		if err != nil {
			self.Close(fmt.Sprintf("write msg to remote error:%s", err))
		}
	} else {
		logger.Error(fmt.Sprintf("socket is close endPoint:%s", self.endPoint))
	}
}

func (self SocketChannel) IsConnect() bool {
	return self.cid != -1
}

func (e *SocketChannel) Close(s string) {
	logger.Info(s)
	e.con.Close()
	e.cid = -1
}
