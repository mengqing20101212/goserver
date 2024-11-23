package server

import (
	"bytes"
	"db"
	"fmt"
	"gameProject/common/utils"
	"logger"
	"net"
	"os"
	"time"
)

// DefaultInputLen defines the default length (5120 bytes) for the buffer used in socket communication.
const DefaultInputLen = 1024 * 5

// DefaultMaxConnectLen specifies the default maximum number of connections that can be managed.
const DefaultMaxConnectLen = 1024

var GeneralCodec = new(PackageFactory) //全局的编码解码器
type ServerInterface interface {
	CreateNewClient(channel *SocketChannel) NetClientInterface
	SetServerPort(port int)
	Start(serverType, serverId, runModule string, startServerCallBack func())
	Stop()
}

// TCP 服务端连接
type Server struct {
	port  int    //端口
	proto string //协议类型 目前是 只支持 tcp

	filterChain   *FilterChain     // filterChain represents a sequence of filters to process packages in the server.
	listener      *net.TCPListener // listener holds the TCP listener for accepting incoming connections.
	ConnectManger *ConnectManger   // ConnectManger handles the management of active socket connections for the server.
	connectNum    uint16           // connectNum indicates the number of active connections managed by the server.
}
type ServerStatusEnum int32

const (
	OPEN   ServerStatusEnum = iota
	WHILTE ServerStatusEnum = 1
	CLOSE  ServerStatusEnum = 2
)

type ServerNodeStatus struct {
	ServerId   string `json:"serverId"`
	ServerType string `json:"serverType"`
	Addr       string `json:"addr"`
	Status     int    `json:"status"`
	Load       int    `json:"load"`
	RunModule  string `json:"runModule"`
}

var log *logger.Logger

// NewServer initializes and returns a new Server instance with the specified port.
func NewServer(port int) (server Server) {
	server = Server{
		port:          port,
		proto:         "tcp",
		filterChain:   &FilterChain{},
		ConnectManger: NewConnectManger(DefaultMaxConnectLen),
		connectNum:    1,
	}
	return server
}

func (this *Server) CreateNewClient(channel *SocketChannel) NetClientInterface {
	return NewNetClient(channel)
}

func (this *Server) SetServerPort(port int) {
	this.port = port
}

// Start initializes the server, sets up the filter chain, starts listening on the specified port,
// and accepts incoming TCP connections in a loop.
func (self *Server) Start(serverType, serverId, runModule string, startServerCallBack func()) {
	if log == nil {
		log = logger.NetLogger
	}
	log.Info("[Server] start begin")
	self.filterChain.AddFilter(&Filter{})
	self.filterChain.AddFilter(&IpFilter{})
	ipaddr := net.TCPAddr{Port: self.port}
	DefaultInitHandler()
	lis, err := net.ListenTCP("tcp", &ipaddr)
	if err != nil {
		log.Error(fmt.Sprintf("[Server] start ListenTCP error:%s", err))
		os.Exit(10)
	}
	self.listener = lis
	if startServerCallBack != nil {
		startServerCallBack()
	}
	CreateServerStatus(self, serverType, serverId, runModule)
	log.Info(fmt.Sprintf("[Server] start end listener port:%d", self.port))
	for {
		con, err := self.listener.Accept()
		if err != nil {
			log.Error(fmt.Sprintf("[Server] Accept error: %s", err))
			continue
		}
		go self.OnAccept(con, self.connectNum)
		self.connectNum++
	}

}

// OnAccept handles the acceptance of a new TCP connection, initializes a new SocketChannel, and adds it to the ConnectManager.
// cid socketId
func (self *Server) OnAccept(con net.Conn, cid uint16) {
	tcpConn, ok := con.(*net.TCPConn)
	if !ok {
		// 处理类型转换失败的情况
		log.Error(fmt.Sprintf("OnAccept conn tcpConn, ok := con.(*net.TCPConn) error"))
		return
	}
	addr := tcpConn.RemoteAddr()
	sc := SocketChannel{
		endPoint: addr,
		cid:      cid, //socketId
		con:      tcpConn,
		inputMsg: utils.NewByteBufferByBuf(bytes.NewBuffer(make([]byte, DefaultInputLen))),
	}
	netClient := ServerInterface(self).CreateNewClient(&sc)
	//设置socket读写队列大小
	err := netClient.GetSocketChannel().con.SetWriteBuffer(DefaultInputLen)
	if err != nil {
		log.Error(fmt.Sprintf("SetWriteBuffer error:%s DefaultInputLen:%d", err, DefaultInputLen))
		return
	}
	err = netClient.GetSocketChannel().con.SetReadBuffer(DefaultInputLen)
	if err != nil {
		log.Error(fmt.Sprintf("setReadBuffer error:%s DefaultInputLen:%d", err, DefaultInputLen))
		return
	}
	sc.inputMsg.GetBuffer().Reset()
	self.ConnectManger.AddConn(netClient, self)
}

func (this *Server) Stop() {
	this.listener.Close()

}

type SocketChannel struct {
	endPoint net.Addr         // endPoint represents the network address of the remote connection for the SocketChannel.
	socketIp string           // socketIp is the IP address associated with the socket connection.
	cid      uint16           // cid socket channel identifier
	con      *net.TCPConn     // con represents the TCP connection associated with the SocketChannel.
	inputMsg utils.ByteBuffer // inputMsg stores the incoming sendMessage data as a ByteBuffer.
}

func (this *SocketChannel) GetCid() uint16 {
	return this.cid
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
		log.Error(fmt.Sprintf("socket is close endPoint:%s", self.endPoint))
	}
}

func (self *SocketChannel) IsConnect() bool {
	return self.cid != 0
}

func (e *SocketChannel) Close(s string) {
	log.Info(s)
	e.con.Close()
	e.cid = 0
}

var ServerStatus ServerNodeStatus

func UpdateServerNodeStatus(status ServerStatusEnum) {
	ServerStatus.Status = int(status)
	jsonStr := utils.ToJsonStr(ServerStatus)
	key := db.RedisKeys(db.GameServerStatusKeyEnum, ServerStatus.ServerType, ServerStatus.ServerId)
	set, err := db.RedisSet(key, *jsonStr, 10*time.Second)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("set redis key:%s value:%s", key, set))
}

func CreateServerStatus(server *Server, serverType, serverId, runModule string) {
	ServerStatus = ServerNodeStatus{
		Addr:       server.listener.Addr().String(),
		ServerType: serverType,
		ServerId:   serverId,
		Load:       0,
		RunModule:  runModule,
	}
	UpdateServerNodeStatus(OPEN)
}
