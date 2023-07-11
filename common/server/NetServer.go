package server

type Server struct {
	port        int
	proto       string
	codecsProto *CodeProto
}

type EndPoint struct {
	ip   string
	port uint16
}

type SocketChannel struct {
	msg      Package
	endPoint EndPoint
}
