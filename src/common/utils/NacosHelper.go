package utils

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// github.com/nacos-group/nacos-sdk-go/v2/common/constant
// http://139.224.80.204:8848/
var ip = "139.224.80.204"
var NameClient *naming_client.INamingClient
var ConfigClient *config_client.IConfigClient

func InitNacos(serverId, serverType string, parserConfigCallback func(string)) {
	//create ServerConfig
	var sc = []constant.ServerConfig{
		*constant.NewServerConfig(ip, 8848),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		//constant.WithNamespaceId("public"),
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
	ConfigClient = &configClient
	config := initConfig(serverId)
	if parserConfigCallback != nil {
		parserConfigCallback(*config)
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

	NameClient = &client

	success, err := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        8848,
		ServiceName: serverId,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		ClusterName: "cluster-a", // 默认值DEFAULT
		GroupName:   serverType,  // 默认值DEFAULT_GROUP
	})
	if err != nil {
		panic("RegisterInstance error")
	}
	if !success {
		panic("RegisterInstance error")
	}
}

func initConfig(serverId string) *string {
	cfg, err := (*ConfigClient).GetConfig(vo.ConfigParam{
		DataId: "gameConfig",
		Group:  "DEFAULT_GROUP",
	})
	if err != nil {
		panic(err)
	}
	println(cfg)
	return &cfg
	//common.Context.ParserConfig(cfg)

}

func CloseNacos() {
	_, _ = (*NameClient).DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        8848,
		ServiceName: "common.Context.Config.ServerId",
		Ephemeral:   true,
		Cluster:     "cluster-a", // 默认值DEFAULT
		GroupName:   "1231",      // 默认值DEFAULT_GROUP
	})

}
