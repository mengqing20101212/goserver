package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ServerConfig struct {
	Id            int32  `json:"Id"`            //服务器id
	Name          string `json:"Name"`          //服务器名称
	Servertype    int32  `json:"Servertype"`    //类型
	Port          int32  `json:"Port"`          //监听端口
	Groupid       int32  `json:"Groupid"`       //组id
	Groupname     string `json:"Groupname"`     //服务器组名称
	Mysqlip       string `json:"Mysqlip"`       //数据库ip
	Mysqlport     int32  `json:"Mysqlport"`     //数据库端口
	Mysqlusername string `json:"Mysqlusername"` //数据库用户名
	Mysqldbname   string `json:"Mysqldbname"`   //数据库
	Mysqlpassword string `json:"Mysqlpassword"` //数据库密码
	Redisip       string `json:"Redisip"`       //redis ip
	Redisport     int32  `json:"Redisport"`     //redis 端口
	Redispassword string `json:"Redispassword"` //redis 密码
	Redisusername string `json:"Redisusername"` //redis 用户名
	Maxconnectnum int32  `json:"Maxconnectnum"` //最大连接数
	Inputsize     string `json:"Inputsize"`     //socket连接上行数据buffer大小
	Outputsize    string `json:"Outputsize"`    //socket连接下行数据buffer大小
	Zoneid        int32  `json:"Zoneid"`        //战区id
	Partid        int32  `json:"Partid"`        //小区id

	//***** 自定义代码区 filed begin ****

	//***** 自定义代码区 filed end ****
}

type ServerConfigPtr struct {
	configList []*ServerConfig
	configMap  map[int32]*ServerConfig
}

func (self *ServerConfigPtr) GetConfigList() []*ServerConfig {
	return self.configList
}

func (self *ServerConfigPtr) GetServerConfig(id int32) *ServerConfig {
	ptr := self.configMap[id]
	if ptr == nil {
		log.Error(fmt.Sprintf("not found ServerConfig id:%d", id))
	}
	return ptr
}
func (self *ServerConfigPtr) GetServerConfigMap() *map[int32]*ServerConfig {
	return &self.configMap
}
func (self *ServerConfigPtr) setConfigList(cfgList []*ServerConfig) {
	self.configList = cfgList
}

func (self *ServerConfigPtr) setConfigMap(configMap map[int32]*ServerConfig) {
	self.configMap = configMap
}

type ServerConfigSwitch struct {
	config1  ServerConfigPtr
	config2  ServerConfigPtr
	isSwitch bool
	lock     sync.RWMutex
}

var serverConfigSwitch = ServerConfigSwitch{}

func LoadServerConfig(fileUrl string) bool {
	serverConfigSwitch.lock.Lock()
	defer serverConfigSwitch.lock.Unlock()
	fs, err := os.OpenFile(filepath.Join(fileUrl, "server.txt"), os.O_RDWR, 0755)
	if err != nil {
		log.Error(fmt.Sprintf("load fileUrl:%s error:%s", fileUrl, err))
		return false
	}
	defer fs.Close()
	scanner := bufio.NewScanner(fs)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if i := strings.Index(string(data), "\t\r\n"); i >= 0 {
			return i + 3, data[0:i], nil
		}
		// 如果数据流结束，并且剩余的数据不为空，则将剩余的数据作为一个 token 返回
		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}
		// 数据不足以分割一个完整的 token，继续读取更多数据
		return 0, nil, nil
	})
	list := make([]*ServerConfig, 0)
	cfgMap := make(map[int32]*ServerConfig, 0)
	for scanner.Scan() {
		line := scanner.Text()
		var serverConfig ServerConfig
		err = json.Unmarshal([]byte(line), &serverConfig)
		if err != nil {
			panic(fmt.Sprintf("load file:%s, line:%s Unmarshal json error:%s", fileUrl, line, err))
		}
		list = append(list, &serverConfig)
		cfgMap[serverConfig.Id] = &serverConfig
	}
	if serverConfigSwitch.isSwitch {
		serverConfigSwitch.config2.setConfigList(list)
		serverConfigSwitch.config2.setConfigMap(cfgMap)
		serverConfigSwitch.afterLoad(&serverConfigSwitch.config2)
	} else {
		serverConfigSwitch.config1.setConfigList(list)
		serverConfigSwitch.config1.setConfigMap(cfgMap)
		serverConfigSwitch.afterLoad(&serverConfigSwitch.config1)
	}
	serverConfigSwitch.isSwitch = !serverConfigSwitch.isSwitch
	log.Info(fmt.Sprintf("load file:%s, success", fs.Name()))
	return true
}

func GetServerConfigPtr() *ServerConfigPtr {
	serverConfigSwitch.lock.RLock()
	defer serverConfigSwitch.lock.RUnlock()
	if serverConfigSwitch.isSwitch {
		return &serverConfigSwitch.config1
	} else {
		return &serverConfigSwitch.config2
	}
}

//***** 自定义代码区 begin ****

func (self *ServerConfigSwitch) afterLoad(config *ServerConfigPtr) {
}

//***** 自定义代码区 end ****
