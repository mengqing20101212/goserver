package server

import (
	"fmt"
	"goserver/common/logger"
)

type FilterInterface interface {
	DoFilter(msg *SocketChannel) (success bool)
}

type Filter struct {
}

func (self Filter) DoFilter(msg *SocketChannel) bool {
	if logger.IsDebug() {
		logger.Debug(" test default filter")
	}
	return true
}

type IpFilter struct {
	ipMap map[string]bool
}

func (self IpFilter) DoFilter(channel *SocketChannel) bool {
	black := self.ipMap[channel.endPoint.ip]
	if black {
		logger.Error(fmt.Sprintf("IpFilter check ip is black channel:%s", channel.String()))
	}
	return black
}

type FilterChain struct {
	filterList []*FilterInterface
}

func (self FilterChain) Filter(channel *SocketChannel) (success bool) {
	for _, filter := range self.filterList {
		doNext := (*filter).DoFilter(channel)
		if !doNext {
			logger.Info(fmt.Sprintf("the channel filter no pass, %s", channel.String()))
			return false
		}
	}
	return true
}
func (self FilterChain) AddFilter(filter *FilterInterface) {
	self.filterList = append(self.filterList, filter)
}
