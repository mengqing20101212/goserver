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

type ActivitypassawardConfig struct {
	Id             int32 `json:"Id"`             //索引id
	Scheduling     int32 `json:"Scheduling"`     //活动排期
	Level          int32 `json:"Level"`          //等级
	Paylevel       int32 `json:"Paylevel"`       //购买等级所需彩钻
	Rechargeshopid int32 `json:"Rechargeshopid"` //商品id
	Score          int32 `json:"Score"`          //需要积分
	Freegift       int32 `json:"Freegift"`       //免费奖励掉落
	Freegiftshow   int32 `json:"Freegiftshow"`   //免费奖励展示
	Paygift        int32 `json:"Paygift"`        //付费奖励掉落
	Paygiftshow    int32 `json:"Paygiftshow"`    //付费奖励展示
	Redirectionid  int32 `json:"Redirectionid"`  //奖励是否突出显示

	//***** 自定义代码区 filed begin ****

	//***** 自定义代码区 filed end ****
}

type ActivitypassawardConfigPtr struct {
	configList []*ActivitypassawardConfig
	configMap  map[int32]*ActivitypassawardConfig
}

func (self *ActivitypassawardConfigPtr) GetConfigList() []*ActivitypassawardConfig {
	return self.configList
}

func (self *ActivitypassawardConfigPtr) GetActivitypassawardConfig(id int32) *ActivitypassawardConfig {
	ptr := self.configMap[id]
	if ptr == nil {
		logger.SystemLogger.Error(fmt.Sprintf("not found ActivitypassawardConfig id:%d", id))
	}
	return ptr
}
func (self *ActivitypassawardConfigPtr) GetActivitypassawardConfigMap() *map[int32]*ActivitypassawardConfig {
	return &self.configMap
}
func (self *ActivitypassawardConfigPtr) setConfigList(cfgList []*ActivitypassawardConfig) {
	self.configList = cfgList
}

func (self *ActivitypassawardConfigPtr) setConfigMap(configMap map[int32]*ActivitypassawardConfig) {
	self.configMap = configMap
}

type ActivitypassawardConfigSwitch struct {
	config1  ActivitypassawardConfigPtr
	config2  ActivitypassawardConfigPtr
	isSwitch bool
	lock     sync.RWMutex
}

var activitypassawardConfigSwitch = ActivitypassawardConfigSwitch{}

func LoadActivitypassawardConfig(fileUrl string) bool {
	activitypassawardConfigSwitch.lock.Lock()
	defer activitypassawardConfigSwitch.lock.Unlock()
	fs, err := os.OpenFile(filepath.Join(fileUrl, "activitypassaward.txt"), os.O_RDWR, 0755)
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
	list := make([]*ActivitypassawardConfig, 0)
	cfgMap := make(map[int32]*ActivitypassawardConfig, 0)
	for scanner.Scan() {
		line := scanner.Text()
		var activitypassawardConfig ActivitypassawardConfig
		err = json.Unmarshal([]byte(line), &activitypassawardConfig)
		if err != nil {
			panic(fmt.Sprintf("load file:%s, line:%s Unmarshal json error:%s", fileUrl, line, err))
		}
		list = append(list, &activitypassawardConfig)
		cfgMap[activitypassawardConfig.Id] = &activitypassawardConfig
	}
	if activitypassawardConfigSwitch.isSwitch {
		activitypassawardConfigSwitch.config2.setConfigList(list)
		activitypassawardConfigSwitch.config2.setConfigMap(cfgMap)
		activitypassawardConfigSwitch.afterLoad(&activitypassawardConfigSwitch.config2)
	} else {
		activitypassawardConfigSwitch.config1.setConfigList(list)
		activitypassawardConfigSwitch.config1.setConfigMap(cfgMap)
		activitypassawardConfigSwitch.afterLoad(&activitypassawardConfigSwitch.config1)
	}
	activitypassawardConfigSwitch.isSwitch = !activitypassawardConfigSwitch.isSwitch
	logger.SystemLogger.Info(fmt.Sprintf("load file:%s, success", fs.Name()))
	return true
}

func GetActivitypassawardConfigPtr() *ActivitypassawardConfigPtr {
	activitypassawardConfigSwitch.lock.RLock()
	defer activitypassawardConfigSwitch.lock.RUnlock()
	if activitypassawardConfigSwitch.isSwitch {
		return &activitypassawardConfigSwitch.config1
	} else {
		return &activitypassawardConfigSwitch.config2
	}
}

//***** 自定义代码区 begin ****

func (self *ActivitypassawardConfigSwitch) afterLoad(config *ActivitypassawardConfigPtr) {
}

//***** 自定义代码区 end ****
