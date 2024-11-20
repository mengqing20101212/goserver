package db

import "fmt"

type RedisKey struct {
	key  string //key 字符串
	desc string //key 描述
}

var (
	GameServerStatusKeyEnum RedisKey = RedisKey{key: "GameServerStatus:%s:%s", desc: "服务器状态GameServerStatus:ServerType:ServerId "}
)

func RedisKeys(serverKeys RedisKey, param ...string) string {
	return fmt.Sprint(serverKeys.key, param)
}

var (
	PlayerServerIdMap RedisKey = RedisKey{key: "PlayerServerId:%s", desc: "玩家当前在那个节点服务器映射"}
)
