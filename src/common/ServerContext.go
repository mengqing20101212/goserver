package common

import (
	"common/utils"
	"gopkg.in/yaml.v3"
	"logger"
	"server"
)

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

func ParserConfig(cfg string) {
	err := yaml.Unmarshal([]byte(cfg), &Context.Config)
	if err != nil {
		panic("parser config error " + err.Error())
		return
	}
	log.Info(cfg)
}

var Context = new(ServerContext)
var log *logger.Logger

func InitContext(logDir, serverId, env string, serverType ServerType) {
	Context.Config.LogDir = logDir
	logger.InitType(logDir)
	log = logger.SystemLogger
	utils.InitNacos(serverId, serverType.String(), env, ParserConfig)
	utils.RegisterNewServerCallBack(serverType.String(), func(serverType string) {
		log.Info("serverType: " + serverType)
	})
}
