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

type ActivityInfoConfig struct {
	Id                   int32  `json:"Id"`                   //编号
	Name                 string `json:"Name"`                 //功能名称
	Opentype             int32  `json:"Opentype"`             //开启类型
	Scheduling           int32  `json:"Scheduling"`           //活动排期
	Openpara1            int32  `json:"Openpara1"`            //开启参数1
	Openpara2            int32  `json:"Openpara2"`            //开启参数2
	Openpara3            int32  `json:"Openpara3"`            //开启参数3
	Timetype             int32  `json:"Timetype"`             //时间类型
	Starttime            string `json:"Starttime"`            //开始时间
	Endtime              string `json:"Endtime"`              //结束时间
	Specialendtime       string `json:"Specialendtime"`       //特殊时间
	Freshtime            string `json:"Freshtime"`            //刷新时间
	Closetime            int32  `json:"Closetime"`            //关闭时间
	Closeactivity        int32  `json:"Closeactivity"`        //任务领取完是否关闭活动
	Openserviceactivity  int32  `json:"Openserviceactivity"`  //开服区间
	Integraltype         int32  `json:"Integraltype"`         //积分类型
	Integralstage        string `json:"Integralstage"`        //阶段积分
	Integralreward       string `json:"Integralreward"`       //积分奖励（掉落表id）
	Integralrewardshow   string `json:"Integralrewardshow"`   //积分奖励（前端）
	Title                string `json:"Title"`                //活动标题
	Picture              string `json:"Picture"`              //立绘
	Description          string `json:"Description"`          //立绘描述
	Para1                string `json:"Para1"`                //功能参数1
	Para2                string `json:"Para2"`                //功能参数2
	Para3                string `json:"Para3"`                //功能参数3
	Mailtemplateid       int32  `json:"Mailtemplateid"`       //邮件模板id
	Exchangeresources    string `json:"Exchangeresources"`    //活动剩余道具转换资源
	Entertype            int32  `json:"Entertype"`            //活动入口类型
	Sort                 int32  `json:"Sort"`                 //排序
	Des                  int32  `json:"Des"`                  //活动描述
	Destime              string `json:"Destime"`              //活动时间描述
	Timedown             int32  `json:"Timedown"`             //是否显示倒计时
	Despic               string `json:"Despic"`               //活动描述
	Rechargeid           string `json:"Rechargeid"`           //商品id
	Iactivitytype        int32  `json:"Iactivitytype"`        //活动类型
	Topid                int32  `json:"Topid"`                //topid
	Noshow               int32  `json:"Noshow"`               //是否不显示在活动栏
	Displayfunctiontype  int32  `json:"Displayfunctiontype"`  //功能显示解锁类型
	Displayfunctionparam string `json:"Displayfunctionparam"` //解锁类型参数

	//***** 自定义代码区 filed begin ****

	//***** 自定义代码区 filed end ****
}

type ActivityInfoConfigPtr struct {
	configList []*ActivityInfoConfig
	configMap  map[int32]*ActivityInfoConfig
}

func (self *ActivityInfoConfigPtr) GetConfigList() []*ActivityInfoConfig {
	return self.configList
}

func (self *ActivityInfoConfigPtr) GetActivityInfoConfig(id int32) *ActivityInfoConfig {
	ptr := self.configMap[id]
	if ptr == nil {
		log.Error(fmt.Sprintf("not found ActivityInfoConfig id:%d", id))
	}
	return ptr
}
func (self *ActivityInfoConfigPtr) GetActivityInfoConfigMap() *map[int32]*ActivityInfoConfig {
	return &self.configMap
}
func (self *ActivityInfoConfigPtr) setConfigList(cfgList []*ActivityInfoConfig) {
	self.configList = cfgList
}

func (self *ActivityInfoConfigPtr) setConfigMap(configMap map[int32]*ActivityInfoConfig) {
	self.configMap = configMap
}

type ActivityInfoConfigSwitch struct {
	config1  ActivityInfoConfigPtr
	config2  ActivityInfoConfigPtr
	isSwitch bool
	lock     sync.RWMutex
}

var activityinfoConfigSwitch = ActivityInfoConfigSwitch{}

func LoadActivityInfoConfig(fileUrl string) bool {
	activityinfoConfigSwitch.lock.Lock()
	defer activityinfoConfigSwitch.lock.Unlock()
	fs, err := os.OpenFile(filepath.Join(fileUrl, "activityInfo.txt"), os.O_RDWR, 0755)
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
	list := make([]*ActivityInfoConfig, 0)
	cfgMap := make(map[int32]*ActivityInfoConfig, 0)
	for scanner.Scan() {
		line := scanner.Text()
		var activityinfoConfig ActivityInfoConfig
		err = json.Unmarshal([]byte(line), &activityinfoConfig)
		if err != nil {
			panic(fmt.Sprintf("load file:%s, line:%s Unmarshal json error:%s", fileUrl, line, err))
		}
		list = append(list, &activityinfoConfig)
		cfgMap[activityinfoConfig.Id] = &activityinfoConfig
	}
	if activityinfoConfigSwitch.isSwitch {
		activityinfoConfigSwitch.config2.setConfigList(list)
		activityinfoConfigSwitch.config2.setConfigMap(cfgMap)
		activityinfoConfigSwitch.afterLoad(&activityinfoConfigSwitch.config2)
	} else {
		activityinfoConfigSwitch.config1.setConfigList(list)
		activityinfoConfigSwitch.config1.setConfigMap(cfgMap)
		activityinfoConfigSwitch.afterLoad(&activityinfoConfigSwitch.config1)
	}
	activityinfoConfigSwitch.isSwitch = !activityinfoConfigSwitch.isSwitch
	log.Info(fmt.Sprintf("load file:%s, success", fs.Name()))
	return true
}

func GetActivityInfoConfigPtr() *ActivityInfoConfigPtr {
	activityinfoConfigSwitch.lock.RLock()
	defer activityinfoConfigSwitch.lock.RUnlock()
	if activityinfoConfigSwitch.isSwitch {
		return &activityinfoConfigSwitch.config1
	} else {
		return &activityinfoConfigSwitch.config2
	}
}

//***** 自定义代码区 begin ****

func (self *ActivityInfoConfigSwitch) afterLoad(config *ActivityInfoConfigPtr) {
}

//***** 自定义代码区 end ****
