package common

import "sync"

var ServerNodeMap = make(map[string]*ServerNode)
var ServerNodeMapLock = &sync.RWMutex{}

func RegisterServerNode(serverNode *ServerNode) {
	ServerNodeMapLock.Lock()
	defer ServerNodeMapLock.Unlock()
	ServerNodeMap[serverNode.ServerId] = serverNode
}
func UnRegisterServerNode(serverId string) {
	ServerNodeMapLock.Lock()
	defer ServerNodeMapLock.Unlock()
	delete(ServerNodeMap, serverId)
}

func GetServerNode(serverId string) *ServerNode {
	ServerNodeMapLock.RLock()
	defer ServerNodeMapLock.RUnlock()
	return ServerNodeMap[serverId]
}

func GetServerNodeMapByServerType(serverType ServerType) map[string]*ServerNode {
	ServerNodeMapLock.RLock()
	defer ServerNodeMapLock.RUnlock()
	resultMap := make(map[string]*ServerNode)
	for k, v := range ServerNodeMap {
		if v.ServerType == serverType {
			resultMap[k] = v
		}
	}
	return resultMap
}
