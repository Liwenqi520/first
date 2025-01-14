package helper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Liwenqi520/errorx"
	"github.com/Liwenqi520/redisx"
)

// 唯一键
type UUID struct{}

var (
	_redisClient redisx.Client
)

// Init redis初始化连接池
func (UUID) Init(config redisx.Config) {
	_redisClient = redisx.New(config)
}

// NewID 获取唯一主键
func (UUID) GetUniqueKey() (ID string, err error) {
	ID = ""
	times := 0
	// 获取开一个redis客户端连接
	ctx := context.Background()
	for {
		// 最大尝试次数
		times++
		if times > 100 {
			fmt.Println("GetUniqueKey err,", "GetUniqueKey,retry:"+strconv.Itoa(times))
			return "", errorx.New(100000, "获取唯一键失败")
		}
		// 获取纳秒
		now := time.Now()
		ID = strconv.FormatInt(now.UnixNano()/1000, 10) // 取微妙,并转为字符串
		// 通过redis setnx 保证id是唯一的
		key := "get_primary_key_id" + ID
		// 判断id是否唯一
		success, err := _redisClient.SetNX(ctx, key, 1, time.Second*1).Result()
		if !success || err != nil {
			time.Sleep(5 * time.Microsecond) // 休眠2微妙，避免最大尝试次数内获得的是同一个微妙
			continue
		}
		break
	}
	return
}
