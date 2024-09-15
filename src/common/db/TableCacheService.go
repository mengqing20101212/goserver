package db

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

const tickSaveTimer = 10 * time.Second

type TableCacheService struct {
	cacheMap map[unsafe.Pointer]func()
	lock     sync.RWMutex
	initFlag bool
}

var cacheService = TableCacheService{}

func InitCacheService() {
	cacheService.lock = sync.RWMutex{}
	cacheService.cacheMap = make(map[unsafe.Pointer]func(), 1024)
	cacheService.initFlag = true
	log.Info(fmt.Sprintf("init TableCacheService caches size:%d", len(cacheService.cacheMap)))
	go tickSaveData(&cacheService)
}

func GetCacheService() *TableCacheService {
	return &cacheService
}

func (self *TableCacheService) AddCacheFunc(ptr unsafe.Pointer, cacheFunc func()) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.cacheMap[ptr] = cacheFunc
	log.Info(fmt.Sprintf("add cache map ptr:%p, size:%d", ptr, len(self.cacheMap)))
}
func (self *TableCacheService) DelCacheFunc(ptr unsafe.Pointer) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.cacheMap, ptr)
	log.Info(fmt.Sprintf("delete ptr:%p from cacheMap size:%d", ptr, len(self.cacheMap)))
}

func (self *TableCacheService) ClearAllFunc() {
	self.lock.Lock()
	defer self.lock.Unlock()
	cacheService.cacheMap = make(map[unsafe.Pointer]func(), 1024)
}

func tickSaveData(service *TableCacheService) {
	if !service.initFlag {
		log.Error(fmt.Sprintf("service not initFlag, TableCacheService:%s", service))
		return
	}
	for {
		time.Sleep(tickSaveTimer)
		log.Info(fmt.Sprintf("tickSaveData: size:%d", len(service.cacheMap)))
		service.lock.RLock()
		for _, f := range service.cacheMap {
			f()
		}
		service.lock.RUnlock()
	}
}
