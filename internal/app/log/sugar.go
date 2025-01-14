package log

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

//==================

// 快捷方式
type H map[string]interface{}

// 处理格式化
func handleH(args []interface{}) []interface{} {
	for i, arg := range args {
		m, ok := arg.(H)
		if ok {
			strArg := "{"
			index := 1
			for k, v := range m {
				if index < len(m) {
					strArg += fmt.Sprint(k, ": ", v, ",")
				} else {
					strArg += fmt.Sprint(k, ": ", v)
				}
				index++
			}
			strArg = strArg + "}"
			args[i] = strArg
		}
	}
	return args
}

//================== with ctx =====================

// Operate类型-变更信息
type ChangeInfo struct {
	CompanyID  string      `json:"company_id"`  //企业ID
	OperatorID string      `json:"operator_id"` //操作者ID
	Data       string      `json:"data"`        //操作的内容
	Before     interface{} `json:"before"`      //修改前
	After      interface{} `json:"after"`       //修改后
	DataName   string      `json:"data_name"`   //更改的数据名称
	IP         string      `json:"ip"`          //ip
	ProjectID  string      `json:"project_id"`  //系统id
}

// 用户操作日志 changeInfo只取第一个
func User(ctx context.Context, logType string, changeInfo /*json数据保存*/ ...ChangeInfo) {
	if ctx == nil {
		ctx = context.Background()
	}
	var sugar = logger.sugaredLogger
	if id, ok := FromTraceIDContext(ctx); ok {
		sugar = sugar.With("trace_id", id)
	}
	sugar = sugar.With("user_id", changeInfo[0].OperatorID)
	if id, ok := FromUserIDContext(ctx); ok {
		sugar = sugar.With("user_id", id)
	}
	//日志类型
	if logType != "" {
		sugar = sugar.With("log_type", logType)
	}
	//用户日志标记
	sugar = sugar.With("tag", "user")
	sugar = sugar.With("read_status", "0")
	sugar = sugar.With("time_unix", time.Now().Unix())

	var dataStr string
	if len(changeInfo) > 0 {
		dataBytes, err := json.Marshal(changeInfo[0])
		if err != nil {
			dataStr = ""
		} else {
			dataStr = string(dataBytes)
		}
	}

	sugar = sugar.With("data", dataStr)

	//
	sugar.Info([]interface{}{}...)
}

func Info(ctx context.Context, logType string, args ...interface{}) {
	if logger == nil {
		return
	}
	args = handleH(args)
	startSpan(ctx, logType).Info(args...)
}

func Error(ctx context.Context, logType string, args ...interface{}) {
	if logger == nil {
		return
	}
	args = handleH(args)
	startSpan(ctx, logType).Error(args...)
}

func Fatal(ctx context.Context, logType string, args ...interface{}) {
	if logger == nil {
		return
	}
	args = handleH(args)
	startSpan(ctx, logType).Fatal(args...)
}

func Debug(ctx context.Context, logType string, args ...interface{}) {
	if logger == nil {
		return
	}
	args = handleH(args)
	startSpan(ctx, logType).Debug(args...)
}

func Panic(ctx context.Context, logType string, args ...interface{}) {
	if logger == nil {
		return
	}
	args = handleH(args)
	startSpan(ctx, logType).Panic(args...)
}

func startSpan(ctx context.Context, logType string) *zap.SugaredLogger {
	if ctx == nil {
		ctx = context.Background()
	}
	var sugar = logger.sugaredLogger
	if id, ok := FromTraceIDContext(ctx); ok {
		sugar = sugar.With("trace_id", id)
	}
	if id, ok := FromUserIDContext(ctx); ok {
		sugar = sugar.With("user_id", id)
	}
	//日志类型
	if logType != "" {
		sugar = sugar.With("log_type", logType)
	}
	sugar = sugar.With("time_unix", time.Now().Unix())
	return sugar
}
