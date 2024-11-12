package common

import (
	"common/utils"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"gopkg.in/yaml.v3"
	"logger"
	"server"
	"strconv"
)

type ServerType int

const (
	Game ServerType = iota
	LOGIN
	GATE
	SCENE
	GM
	Unknown
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

	DbConfig    DbConfig    `yaml:"db"`
	RedisConfig RedisConfig `yaml:"redis"`

	//server 相关
	ServerPort int `yaml:"serverPort"`
	ServerType ServerType
	ServerId   string
	//TODO http相关

}

type RedisConfig struct {
	RedisIp       string `yaml:"host"`
	RedisPort     int    `yaml:"port"`
	RedisPassword string `yaml:"password"`
}

type DbConfig struct {
	//数据库相关
	DbIp       string `yaml:"host"`
	DbPort     int    `yaml:"port"`
	DbUser     string `yaml:"userName"`
	DbPassword string `yaml:"passWord"`
	DbName     string `yaml:"dbName"`
}

type ServerContext struct {
	Config ServerConfig
	Server *server.Server
}

// 用于rpc 调用需要的服务器节点信息
type ServerNode struct {
	ServerType ServerType     //服务器类型
	ServerId   string         //服务器id
	ServerPort int            //服务器端口
	data       model.Instance //服务器其他信息 IP 之类的 nacos 信息
}

func (this *ServerNode) String() string {
	return fmt.Sprintln("serverNodeInfo: serverType: ", this.ServerType.String(), " serverId: ", this.ServerId+" ipaddr: ", this.data.Ip, ":", strconv.Itoa(this.ServerPort))
}

func (this *ServerNode) GetIP() string {
	return this.data.Ip
}

func ParserConfig(cfg string) int {
	err := yaml.Unmarshal([]byte(cfg), &Context.Config)
	if err != nil {
		panic("parser config error " + err.Error())
		return 0
	}
	log.Info(cfg)
	return Context.Config.ServerPort
}

func ParserServerNode(isAdd bool, data model.Instance) {
	serverId := data.Metadata["serverId"]
	if isAdd {
		serverNode := new(ServerNode)
		serverNode.ServerId = serverId
		serverNode.ServerPort, _ = strconv.Atoi(data.Metadata["serverPort"])
		serverNode.ServerType = getServerType(data.ServiceName)
		serverNode.data = data
		RegisterServerNode(serverNode)
	} else {
		UnRegisterServerNode(serverId)
	}

}

func getServerType(serviceName string) ServerType {
	switch serviceName {
	case "Game":
		return Game
	case "Login":
		return LOGIN
	case "Gate":
		return GATE
	case "Scene":
		return SCENE
	case "Gm":
		return GM
	default:
		return Unknown
	}
}

var Context = new(ServerContext)
var log *logger.Logger

func InitContext(logDir, serverId, env string, serverType ServerType) {
	Context.Config.LogDir = logDir
	logger.InitType(logDir)
	log = logger.SystemLogger
	utils.InitNacos(serverId, serverType.String(), env, ParserConfig)
	utils.RegisterNewServerCallBack(serverType.String(), ParserServerNode)
}
