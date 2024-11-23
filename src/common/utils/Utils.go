package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

func GetYYYY_MM_DD_HH_mm_ss() string {
	return time.Now().Format("2006_01_01_15_04_05")
}

// 偏移的时间戳
var offsetTimer int64 = 0
var timeLock sync.Mutex

// GetNow 获取服务器当前时间戳 毫秒
func GetNow() int64 {
	return time.Now().UnixMilli() + offsetTimer
}

func SetOffsetTimer(timer int64) {
	timeLock.Lock()
	defer timeLock.Unlock()
	offsetTimer = timer
	fmt.Println("SetOffsetTimer", timer)
}

func ClearOffsetTimer() {
	SetOffsetTimer(0)
}

// 获取服务器当前时间戳 毫秒
func GetNowString() string {
	now := GetNow()
	return time.UnixMilli(now).Format("2006-01-02 15:04:05")
}
func IsSameDay(t1, t2 int64) bool {
	return time.UnixMilli(t1).Day() == time.UnixMilli(t2).Day()
}

// 结构体转换为 byte 数组
func Struct2Bytes(data any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// byte 数组转换为 结构体
func Bytes2Struct(data []byte, result *any) error {
	var buf bytes.Buffer
	buf.Write(data)
	dec := gob.NewDecoder(&buf)
	return dec.Decode(result)
}
func ToJsonStr(obj any) *string {
	// 将结构体转换为JSON格式的字节切片
	jsonBytes, err := json.Marshal(obj)
	str := "{}"
	if err != nil {
		fmt.Println("转换失败:", err)
		return &str
	}
	str = string(jsonBytes)
	return &str
}
func JsonToObj(jsonStr *string, obj any) error {
	err := json.Unmarshal([]byte(*jsonStr), obj)
	if err != nil {
		return err
	}
	return nil
}
