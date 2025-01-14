package log

import (
	"first/internal/defines"
	"context"
	"fmt"
	"time"

	"github.com/Liwenqi520/redisx"
)

type Hooker struct {
	rds redisx.Client
}

func (h *Hooker) Write(p []byte) (n int, err error) {
	logString := string(p)
	r := h.rds
	logKey := defines.LOG
	if _, err := r.LPush(context.Background(), logKey, logString).Result(); err != nil {
		//日志写入失败
		fmt.Println("写入失败")
		return 0, nil
	}
	//检查长度
	lens, err := r.LLen(context.Background(), logKey).Result()
	if err != nil {
		fmt.Println("获取长度失败")
		return 0, nil
	}
	if lens > 1000000 {
		//del
		fmt.Println("日志超过100万")
		if _, err := r.LTrim(context.Background(), logKey, 0, 99).Result(); err != nil {
			//日志写入失败
			fmt.Println("保留100条失败")
			return 0, nil
		}
	}

	//日志时间10分钟
	if _, err := r.Expire(context.Background(), logKey, 600).Result(); err != nil {
		//日志写入失败
		fmt.Println("过期设置失败")
		return 0, nil
	}
	return
}

type HookLog struct {
	Level      string `json:"level"`
	Time       string `json:"time"`
	Caller     string `json:"caller"`
	Msg        string `json:"msg"`
	Version    string `json:"version"`
	ServerName string `json:"server_name"`
	LogType    string `json:"log_type"`
	TraceID    string `json:"trace_id"`
	UserID     string `json:"user_id"`
}

var hk *Hooker

func InitHook(rds redisx.Client) (err error) {
	hk = &Hooker{rds: rds}
	return nil
}

type Logs struct {
	ID         uint      `gorm:"column:id;primary_key;auto_increment;"`                              // id
	ServerName string    `gorm:"column:server_name;size:50;not null;default:'';comment:服务名称;index;"` // 服务名称
	LogType    string    `gorm:"column:log_type;size:50;not null;default:'';comment:类型;index;"`      // 类型
	Level      string    `gorm:"column:level;size:20;not null;default:'';comment:日志级别;index;"`       // 日志级别
	Message    string    `gorm:"column:message;size:1024;not null;default:'';comment:消息;"`           // 消息
	TraceID    string    `gorm:"column:trace_id;size:128;not null;default:'';comment:跟踪id;index;"`   // 跟踪ID
	UserID     string    `gorm:"column:user_id;size:36;not null;default:'';comment:用户id;index;"`     // 用户ID
	Caller     string    `gorm:"column:caller;size:256;not null;default:'';comment:caller信息;"`       // Caller
	Data       string    `gorm:"column:data;type:text;not null;comment:日志数据;"`                       // 日志数据(json)
	Version    string    `gorm:"column:version;index;size:32;not null;default:'';comment:版本号;"`      // 服务版本号
	CreatedAt  time.Time `gorm:"column:created_at;index"`                                            // 创建时间
}

// TableName 表名
func (Logs) TableName() string {
	return "logs"
}
