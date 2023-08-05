package table

import (
	"fmt"
	"goserver/common/db"
	"goserver/common/logger"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

func TestNewSysTable(t *testing.T) {
	logger.Init("../logs", "test2.log")
	db.InitDefaultDataBase("root", "root", "127.0.0.1", "sysweb", 3306)
	db.InitCacheService()
	list := make([]unsafe.Pointer, 100)
	for i := 0; i < 100; i++ {
		sysTable := NewSysTable()
		sysTable.SetStatus(i%2 == 0).SetRole_name(string("test" + strconv.Itoa(i))).SetDescription("meng qing" + strconv.Itoa(i))
		list = append(list, unsafe.Pointer(sysTable))
		fmt.Println("FindOneSysTable new data: ", FindOneSysTable(uint64(i)))
	}
	time.Sleep(10 * time.Second)
	for _, ptr := range list {
		db.GetCacheService().DelCacheFunc(ptr)
	}
	time.Sleep(1024 * time.Second)
}
