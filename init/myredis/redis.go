package myredis

import "github.com/Liwenqi520/redisx"

var (
	RedisClient redisx.Client
)

// Init redis初始化连接池
func RedisInit(config redisx.Config) {
	RedisClient = redisx.New(config)
}

func RedisCli() redisx.Client {
	return RedisClient
}
