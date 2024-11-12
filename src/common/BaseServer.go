package common

type BaseServer interface {
	StartServer(serverId, env string)
	StopServer()
	InitHandler()
	GetServerContext() *ServerContext
}
