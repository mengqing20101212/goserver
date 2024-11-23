package utils

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"logger"
	"strconv"
)

// github.com/nacos-group/nacos-sdk-go/v2/common/constant
// http://139.224.80.204:8848/
var ip = "192.168.161.182"
var serverName = "rpcNodeService"
var NameClientPtr *naming_client.INamingClient
var ConfigClientPtr *config_client.IConfigClient
var log *logger.Logger

// InitNacos 初始化Nacos客户端，包括配置客户端和命名客户端，并注册服务实例。
// serverId: 服务器ID
// serverType: 服务类型
// env: 环境变量
// parserConfigCallback: 配置解析回调函数
func InitNacos(serverId, serverType, env string, parserConfigCallback func(string) int) {
	//create ServerConfig
	if log == nil {
		log = logger.SystemLogger
	}
	log.Info(fmt.Sprintln("InitNacos params serverId:", serverId, "serverType:", serverType, "env:", env))
	var sc = []constant.ServerConfig{
		*constant.NewServerConfig(ip, 8848),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(env),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("../../nacos/log"),
		constant.WithCacheDir("../../nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	ConfigClientPtr = &configClient
	config := initConfig(&serverId, &serverType)
	serverPort := 0
	if parserConfigCallback != nil {
		serverPort = parserConfigCallback(*config)
	}

	// create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	NameClientPtr = &client

	//注册该节点
	serverConfig := make(map[string]string)
	serverConfig["serverId"] = serverId
	serverConfig["serverPort"] = strconv.Itoa(serverPort)
	success, err := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        8848,
		ServiceName: serverName,
		Weight:      100,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		ClusterName: env,        // 默认值DEFAULT
		GroupName:   serverType, // 默认值DEFAULT_GROUP
		Metadata:    serverConfig,
	})
	if err != nil {
		print(err)
		panic("RegisterInstance error")
	}

	//监听节点
	if !success {
		panic("RegisterInstance error")
	}
}
func RegisterNewServerCallBack(serverType string, newServerCallBack func(isAdd bool, instance model.Instance)) {
	if NameClientPtr == nil {
		log.Error("NameClientPtr is nil")
		return
	}

	err := (*NameClientPtr).Subscribe(&vo.SubscribeParam{
		ServiceName: serverName,
		GroupName:   serverType,
		SubscribeCallback: func(services []model.Instance, err error) {
			if err != nil {
				log.Error("订阅回调发生错误: " + err.Error())
				return
			}
			log.Info(fmt.Sprintln(" 上架节点,当前节点数量:", len(services)))
			for _, v := range services {
				log.Info(fmt.Sprintln("服务实例信息: ", v))
				newServerCallBack(v.Enable, v)
			}
		},
	})
	if err != nil {
		log.Error("订阅回调发生错误: " + err.Error())
		return
	}

}

func initConfig(serverId, serverType *string) *string {
	cfg, err := (*ConfigClientPtr).GetConfig(vo.ConfigParam{
		//DataId:  *serverId,
		DataId:  "game1001",
		Group:   *serverType,
		Content: "",
	})
	if err != nil {
		print(err)
		panic(err)
	}
	log.Info(cfg)
	//common.Context.ParserConfig(cfg)
	return &cfg

}

func CloseNacos() {
	_, _ = (*NameClientPtr).DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        8848,
		ServiceName: "common.Context.Config.ServerId",
		Ephemeral:   true,
		Cluster:     "cluster-a", // 默认值DEFAULT
		GroupName:   "1231",      // 默认值DEFAULT_GROUP
	})

}
