package config

import (
	"bufio"
	"common/logger"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ActivityNpcConfig struct {
	Id               int32  `json:"Id"`               //编号
	Config_name      string `json:"Config_name"`      //配置名字辅助列
	Charactermodelid string `json:"Charactermodelid"` //活动名称
	Defaultani       string `json:"Defaultani"`       //默认动作
	Bornshowani      string `json:"Bornshowani"`      //欢迎动作
	Clicktype        string `json:"Clicktype"`        //点击反馈类型(0-不可点，1-通用点击逻辑，2-特殊
	Clickanilist     string `json:"Clickanilist"`     //点击反馈动作组
	Clickcameralist  string `json:"Clickcameralist"`  //点击反馈镜头组
	Movedistance     string `json:"Movedistance"`     //推镜的镜头左右偏移（正是左负是右,左右是角度，上下是距离）
	Cameradistance   string `json:"Cameradistance"`   //推镜的距离配置
	Clicktext        string `json:"Clicktext"`        //点击反馈文本
	Grouptext        string `json:"Grouptext"`        //组合动作的文本
	Activitytype     string `json:"Activitytype"`     //特殊功能类型
	Param_1          string `json:"Param_1"`          //特殊功能参数1
	Param_2          string `json:"Param_2"`          //特殊功能参数2
	Npcgrounpid      string `json:"Npcgrounpid"`      //角色组合id
	Showpriority     string `json:"Showpriority"`     //主城显示优先级
	Decorationid     string `json:"Decorationid"`     //装饰物资源id
	Decorationpoint  string `json:"Decorationpoint"`  //资源物挂点

	//***** 自定义代码区 filed begin ****

	//***** 自定义代码区 filed end ****
}

type ActivityNpcConfigPtr struct {
	configList []*ActivityNpcConfig
	configMap  map[int32]*ActivityNpcConfig
}

func (self *ActivityNpcConfigPtr) GetConfigList() []*ActivityNpcConfig {
	return self.configList
}

func (self *ActivityNpcConfigPtr) GetActivityNpcConfig(id int32) *ActivityNpcConfig {
	ptr := self.configMap[id]
	if ptr == nil {
		logger.SystemLogger.Error(fmt.Sprintf("not found ActivityNpcConfig id:%d", id))
	}
	return ptr
}
func (self *ActivityNpcConfigPtr) GetActivityNpcConfigMap() *map[int32]*ActivityNpcConfig {
	return &self.configMap
}
func (self *ActivityNpcConfigPtr) setConfigList(cfgList []*ActivityNpcConfig) {
	self.configList = cfgList
}

func (self *ActivityNpcConfigPtr) setConfigMap(configMap map[int32]*ActivityNpcConfig) {
	self.configMap = configMap
}

type ActivityNpcConfigSwitch struct {
	config1  ActivityNpcConfigPtr
	config2  ActivityNpcConfigPtr
	isSwitch bool
	lock     sync.RWMutex
}

var activitynpcConfigSwitch = ActivityNpcConfigSwitch{}

func LoadActivityNpcConfig(fileUrl string) bool {
	activitynpcConfigSwitch.lock.Lock()
	defer activitynpcConfigSwitch.lock.Unlock()
	fs, err := os.OpenFile(filepath.Join(fileUrl, "activityNpc.txt"), os.O_RDWR, 0755)
	if err != nil {
		logger.SystemLogger.Error(fmt.Sprintf("load fileUrl:%s error:%s", fileUrl, err))
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
	list := make([]*ActivityNpcConfig, 0)
	cfgMap := make(map[int32]*ActivityNpcConfig, 0)
	for scanner.Scan() {
		line := scanner.Text()
		var activitynpcConfig ActivityNpcConfig
		err = json.Unmarshal([]byte(line), &activitynpcConfig)
		if err != nil {
			panic(fmt.Sprintf("load file:%s, line:%s Unmarshal json error:%s", fileUrl, line, err))
		}
		list = append(list, &activitynpcConfig)
		cfgMap[activitynpcConfig.Id] = &activitynpcConfig
	}
	if activitynpcConfigSwitch.isSwitch {
		activitynpcConfigSwitch.config2.setConfigList(list)
		activitynpcConfigSwitch.config2.setConfigMap(cfgMap)
		activitynpcConfigSwitch.afterLoad(&activitynpcConfigSwitch.config2)
	} else {
		activitynpcConfigSwitch.config1.setConfigList(list)
		activitynpcConfigSwitch.config1.setConfigMap(cfgMap)
		activitynpcConfigSwitch.afterLoad(&activitynpcConfigSwitch.config1)
	}
	activitynpcConfigSwitch.isSwitch = !activitynpcConfigSwitch.isSwitch
	logger.SystemLogger.Info(fmt.Sprintf("load file:%s, success", fs.Name()))
	return true
}

func GetActivityNpcConfigPtr() *ActivityNpcConfigPtr {
	activitynpcConfigSwitch.lock.RLock()
	defer activitynpcConfigSwitch.lock.RUnlock()
	if activitynpcConfigSwitch.isSwitch {
		return &activitynpcConfigSwitch.config1
	} else {
		return &activitynpcConfigSwitch.config2
	}
}

//***** 自定义代码区 begin ****

func (self *ActivityNpcConfigSwitch) afterLoad(config *ActivityNpcConfigPtr) {
}

//***** 自定义代码区 end ****
