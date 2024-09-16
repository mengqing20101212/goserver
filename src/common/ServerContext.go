package common

import "server"

type ServerType int

const (
	Game ServerType = iota
	LOGIN
	GATE
	SCENE
	GM
)

func (s ServerType) String() string {
	switch s {
	case Game:
		return "Game"
	case LOGIN:
		return "LOGIN"
	case GATE:
		return "GATE"
	case SCENE:
		return "SCENE"
	case GM:
		return "GM"
	}
	panic("unknown server type")
}

type ServerConfig struct {
	LogDir string
	//数据库相关
	DbIp       string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string
	//redis相关
	RedisIp       string
	RedisPort     int
	RedisUser     string
	RedisPassword string

	//server 相关
	ServerPort int
	ServerType ServerType
	ServerId   string
	//TODO http相关

}

type ServerContext struct {
	Config *ServerConfig
	Server *server.Server
}

var Context = new(ServerContext)
