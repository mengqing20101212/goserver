package db

import "fmt"

type RedisKey struct {
	key  string //key 字符串
	desc string //key 描述
}

var (
	GameServerStatusKeyEnum RedisKey = RedisKey{key: "GameServerStatus:%s:%s", desc: "服务器状态GameServerStatus:ServerType:ServerId "}
)

func RedisKeys(serverKeys RedisKey, param ...any) string {
	/*if len(param) == 0 {
		return serverKeys.key
	} else if len(param) == 1 {
		return fmt.Sprintf(serverKeys.key, param[0])
	} else if len(param) == 2 {
		return fmt.Sprintf(serverKeys.key, param[0], param[1])
	} else if len(param) == 3 {
		return fmt.Sprintf(serverKeys.key, param[0], param[1], param[2])
	} else if len(param) == 4 {
		return fmt.Sprintf(serverKeys.key, param[0], param[1], param[2], param[3])
	} else if len(param) == 5 {
		return fmt.Sprintf(serverKeys.key, param[0], param[1], param[2], param[3], param[4])
	} else if len(param) == 6 {
		return fmt.Sprintf(serverKeys.key, param[0], param[1], param[2], param[3], param[4], param[5])
	} else if len(param) == 7 {
		return fmt.Sprintf(serverKeys.key, param[0], param[1], param[2], param[3], param[4], param[5], param[6])
	} else if len(param) == 8 {
		return fmt.Sprintf(serverKeys.key, param[0], param[1], param[2], param[3], param[4], param[5], param[6], param[7])
	} else if len(param) == 9 {
		return fmt.Sprintf(serverKeys.key, param[0], param[1], param[2], param[3], param[4], param[5], param[6], param[7], param[8])
	} else if len(param) == 10 {
		return fmt.Sprintf(serverKeys.key, param[0], param[1], param[2], param[3], param[4], param[5], param[6], param[7], param[8], param[9])
	}*/
	return fmt.Sprintf(serverKeys.key, param...)
}

var (
	PlayerServerIdMap RedisKey = RedisKey{key: "PlayerServerId:%s", desc: "玩家当前在那个节点服务器映射"}
)
