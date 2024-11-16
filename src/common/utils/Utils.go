package utils

import (
	"fmt"
	"sync"
	"time"
)

// 偏移的时间戳
var offsetTimer int64 = 0
var timeLock sync.Mutex

// 获取服务器当前时间戳 毫秒
func GetNow() int64 {
	return time.Now().Unix() + offsetTimer
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
