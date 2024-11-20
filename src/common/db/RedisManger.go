package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var redisClient *redis.Client

func InitRedisConnect(ip, port, password, usename string) bool {
	redisClient = redis.NewClient(&redis.Options{
		Addr:           fmt.Sprintf("%s:%s", ip, port),
		Password:       password,
		Username:       usename,
		DB:             0,
		MaxActiveConns: 10,
	})
	var ctx = context.Background()
	res, err := redisClient.Ping(ctx).Result()
	if err != nil {
		DbLogger.Error(fmt.Sprintf("redis connect init error %s", err))
		return false
	}
	DbLogger.Info(fmt.Sprintf("redis connect init success %s", res))
	return true
}
func CloseRedisConnect() {
	if redisClient != nil {
		err := redisClient.Close()
		if err != nil {
			return
		}
	}
}

func RedisGet(key string) (string, error) {
	return redisClient.Get(context.Background(), key).Result()
}
func RedisSet(key, val string, ttl time.Duration) (string, error) {
	return redisClient.Set(context.Background(), key, val, ttl).Result()
}

func GetRedis() *redis.Client {
	return redisClient
}
