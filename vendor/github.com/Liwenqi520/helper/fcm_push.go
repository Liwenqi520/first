package helper

import (
	"context"
	"time"
)

type FcmConfig struct {
	ApiAccessKey string `json:"api_access_key"`
	MessageApi   string `json:"message_api"`
}

func FCMPush(Config FcmConfig, data interface{}) (re interface{}) {
	var timeout time.Duration
	headers := map[string]string{
		"Authorization": "key=" + Config.ApiAccessKey,
		"Content-Type":  "application/json",
	}
	res, _ := PostJSON(context.Background(), Config.MessageApi, data, headers, timeout)
	return res
}
