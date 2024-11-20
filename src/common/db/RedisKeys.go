package db

import "fmt"

type ServerKeys RedisKey
type RedisKey struct {
	key  string //key 字符串
	desc string //key 描述
}

var (
	GameServerStatusKeys ServerKeys = ServerKeys{key: "GameServerStatusKeys:%s", desc: "服务器状态"}
)

func ServerRedisKeys(serverKeys ServerKeys, param ...string) string {
	return fmt.Sprint(serverKeys.key, param)
}

type PlayerKeys RedisKey

var (
	PlayerServerIdMap PlayerKeys = PlayerKeys{key: "PlayerServerId:%s", desc: "玩家当前在那个节点服务器映射"}
)

func PlayerRedisKeys(serverKeys PlayerKeys, param ...string) string {
	return fmt.Sprint(serverKeys.key, param)
}
