package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"goserver/common/logger"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ActivityNpcGroupConfig struct {
	Id             int32  `json:"Id"`             //编号
	Defaultani     string `json:"Defaultani"`     //默认动作
	P1             string `json:"P1"`             //辅助列
	Groupname      string `json:"Groupname"`      //备注是谁
	Clickanilist   string `json:"Clickanilist"`   //点击反馈动作组
	Movedistance   string `json:"Movedistance"`   //推镜的镜头左右偏移（正是左负是右,左右是角度，上下是距离）
	Camerahight    string `json:"Camerahight"`    //注视点高度
	Cameradistance string `json:"Cameradistance"` //推镜的距离配置
	Npcdistance    string `json:"Npcdistance"`    //角色距离
	Textprefabtype string `json:"Textprefabtype"` //文本框的效果类型

	//********filed begin********//

	//********filed end********//
}

type ActivityNpcGroupConfigPtr struct {
	configList []*ActivityNpcGroupConfig
	configMap  map[int32]*ActivityNpcGroupConfig
}

func (self *ActivityNpcGroupConfigPtr) GetConfigList() []*ActivityNpcGroupConfig {
	return self.configList
}

func (self *ActivityNpcGroupConfigPtr) GetActivityNpcGroupConfig(id int32) *ActivityNpcGroupConfig {
	ptr := self.configMap[id]
	if ptr == nil {
		logger.Error(fmt.Sprintf("not found ActivityNpcGroupConfig id:%d", id))
	}
	return ptr
}
func (self *ActivityNpcGroupConfigPtr) GetActivityNpcGroupConfigMap() *map[int32]*ActivityNpcGroupConfig {
	return &self.configMap
}
func (self *ActivityNpcGroupConfigPtr) setConfigList(cfgList []*ActivityNpcGroupConfig) {
	self.configList = cfgList
}

func (self *ActivityNpcGroupConfigPtr) setConfigMap(configMap map[int32]*ActivityNpcGroupConfig) {
	self.configMap = configMap
}

type ActivityNpcGroupConfigSwitch struct {
	config1  ActivityNpcGroupConfigPtr
	config2  ActivityNpcGroupConfigPtr
	isSwitch bool
	lock     sync.RWMutex
}

var activitynpcgroupConfigSwitch = ActivityNpcGroupConfigSwitch{}

func LoadActivityNpcGroupConfig(fileUrl string) bool {
	activitynpcgroupConfigSwitch.lock.Lock()
	defer activitynpcgroupConfigSwitch.lock.Unlock()
	fs, err := os.OpenFile(filepath.Join(fileUrl, "activityNpcGroup.txt"), os.O_RDWR, 0755)
	if err != nil {
		logger.Error(fmt.Sprintf("load fileUrl:%s error:%s", fileUrl, err))
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
	list := make([]*ActivityNpcGroupConfig, 0)
	cfgMap := make(map[int32]*ActivityNpcGroupConfig, 0)
	for scanner.Scan() {
		line := scanner.Text()
		var activitynpcgroupConfig ActivityNpcGroupConfig
		err = json.Unmarshal([]byte(line), &activitynpcgroupConfig)
		if err != nil {
			panic(fmt.Sprintf("load file:%s, line:%s Unmarshal json error:%s", fileUrl, line, err))
		}
		list = append(list, &activitynpcgroupConfig)
		cfgMap[activitynpcgroupConfig.Id] = &activitynpcgroupConfig
	}
	if activitynpcgroupConfigSwitch.isSwitch {
		activitynpcgroupConfigSwitch.config2.setConfigList(list)
		activitynpcgroupConfigSwitch.config2.setConfigMap(cfgMap)
		activitynpcgroupConfigSwitch.afterLoad(&activitynpcgroupConfigSwitch.config2)
	} else {
		activitynpcgroupConfigSwitch.config1.setConfigList(list)
		activitynpcgroupConfigSwitch.config1.setConfigMap(cfgMap)
		activitynpcgroupConfigSwitch.afterLoad(&activitynpcgroupConfigSwitch.config1)
	}
	activitynpcgroupConfigSwitch.isSwitch = !activitynpcgroupConfigSwitch.isSwitch
	logger.Info(fmt.Sprintf("load file:%s, success", fs.Name()))
	return true
}

func GetActivityNpcGroupConfigPtr() *ActivityNpcGroupConfigPtr {
	activitynpcgroupConfigSwitch.lock.RLock()
	defer activitynpcgroupConfigSwitch.lock.RUnlock()
	if activitynpcgroupConfigSwitch.isSwitch {
		return &activitynpcgroupConfigSwitch.config1
	} else {
		return &activitynpcgroupConfigSwitch.config2
	}
}

//*****begin****//

func (self *ActivityNpcGroupConfigSwitch) afterLoad(config *ActivityNpcGroupConfigPtr) {
}

//*****end****//
