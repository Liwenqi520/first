package helper

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
)

var PushIOSTypeAdd = 1000

type AppKeySecretStruct struct {
	AppKey       string `mapstructure:"app_key" json:"app_key"`
	MasterSecret string `mapstructure:"master_secret" json:"master_secret"`
}
type UmengConfig struct {
	HostUrl                 string             `mapstructure:"host_url" json:"host_url"`
	PostUri                 string             `mapstructure:"post_uri" json:"post_uri"`
	PushProductionMode      string             `mapstructure:"push_production_mode" json:"push_production_mode"`             //推送是否是生成模式
	PushAndroidActivityChat string             `mapstructure:"push_android_activity_chat" json:"push_android_activity_chat"` //设置前端收到通知后点击通知跳转的页面
	PushAndroidActivityWeb  string             `mapstructure:"push_android_activity_web" json:"push_android_activity_web"`   //设置前端收到通知后点击通知跳转的页面
	AndroidConfig           AppKeySecretStruct `mapstructure:"android_config" json:"android_config"`
	IOSConfig               AppKeySecretStruct `mapstructure:"ios_config" json:"ios_config"`
}
type UMengAndroid struct {
	AppKey         string          `json:"appkey"`        // 必填项
	Timestamp      string          `json:"timestamp"`     // 必填项
	Type           string          `json:"type"`          // 必填项
	DeviceTokens   string          `json:"device_tokens"` // 选填,用于给特定设备的推送
	ProductionMode string          `json:"production_mode"`
	Payload        *PayloadAndroid `json:"payload"` // 必填项
	Description    string          `json:"description"`
}
type PayloadAndroid struct {
	DisplayType string            `json:"display_type"` // 必填项
	Body        *BodyAndroid      `json:"body"`         // 必填项
	Extra       map[string]string `json:"extra"`
}
type BodyAndroid struct {
	Ticker    string `json:"ticker"`     // 必填项
	Title     string `json:"title"`      // 必填项
	Text      string `json:"text"`       // 必填项
	AfterOpen string `json:"after_open"` // 必填项
	Activity  string `json:"activity"`   // 必填项
}

// android推送
// deviceToken:设备的编号,如果设置deviceToken，则是单播;如果未设置则是全播
func PushAndroid(Config UmengConfig, ticker, title, text string, deviceToken string, extraData map[string]string) (res interface{}, err error) {
	body := &BodyAndroid{}
	// 必填 通知栏提示文字
	body.Ticker = ticker
	// 必填 通知标题
	body.Title = title
	// 必填 通知文字描述
	body.Text = text
	// 打开Android端的Activity
	body.AfterOpen = "go_activity"
	payLoad := &PayloadAndroid{}
	payLoad.DisplayType = "notification"
	payLoad.Body = body
	// 额外携带的信息
	payLoad.Extra = extraData
	messageAndroid := UMengAndroid{}
	messageAndroid.AppKey = Config.AndroidConfig.AppKey
	// 打开聊天
	body.Activity = Config.PushAndroidActivityChat
	if deviceToken == "" {
		// 全播
		messageAndroid.Type = "broadcast"
		// 打开webview
		body.Activity = Config.PushAndroidActivityWeb
	} else {
		// 单播
		messageAndroid.Type = "unicast"
		messageAndroid.DeviceTokens = deviceToken
		// 打开聊天
		body.Activity = Config.PushAndroidActivityChat
	}
	timeInt64 := time.Now().Unix()
	timestamp := strconv.FormatInt(timeInt64, 10)
	messageAndroid.Timestamp = timestamp
	messageAndroid.ProductionMode = Config.PushProductionMode
	messageAndroid.Payload = payLoad
	messageAndroid.Description = title
	postBody, _ := json.Marshal(messageAndroid)
	url := Config.HostUrl + Config.PostUri
	// MD5加密
	sign := Md5("POST" + url + string(postBody) + Config.AndroidConfig.MasterSecret)
	url = url + "?sign=" + sign
	var timeout time.Duration
	res, err = PostJSON(context.Background(), url, messageAndroid, map[string]string{}, timeout)
	return
}

/**
IOS推送必须项:
appkey
"timestamp":"xx",       // 必填 时间戳，10位或者13位均可，时间戳有效期为10分钟
type       //broadcast
"alert": "xx"          // 必填
MasterSecret
"production_mode":"true/false" // 可选 正式/测试模式。测试模式下，只会将消息发给测试设备。
*/
type UMengIOS struct {
	AppKey         string      `json:"appkey"`    // 必填项
	Timestamp      string      `json:"timestamp"` // 必填项
	Type           string      `json:"type"`      // 必填项
	ProductionMode string      `json:"production_mode"`
	Payload        *PayloadIOS `json:"payload"`       // 必填项
	DeviceTokens   string      `json:"device_tokens"` // 选填项
	Description    string      `json:"description"`
}
type PayloadIOS struct {
	Aps   *ApsIOS `json:"aps"`
	PType int     `json:"ptype"` // 1000:咨询,1001:政策法规;1004:法制宣传;1013:新闻
	Purl  string  `json:"purl"`
}
type ApsIOS struct {
	Alert            string `json:"alert"` // 必填项
	ContentAvailable string `json:"content-available"`
}

// IOS推送
// deviceToken:设备的编号,如果设置deviceToken，则是单播;如果未设置则是全播
// pType,purl:额外参数，非必选值
// pType:区分类型
// purl:要打开的超链接的类型
func PushIOS(Config UmengConfig, alert, contentAvailable string, deviceToken string, pType int, purl string) (res interface{}, err error) {
	aps := &ApsIOS{}
	aps.Alert = alert
	aps.ContentAvailable = contentAvailable
	payLoad := &PayloadIOS{}
	payLoad.Aps = aps
	payLoad.PType = pType
	payLoad.Purl = purl
	messageIOS := UMengIOS{}
	messageIOS.Payload = payLoad
	messageIOS.AppKey = Config.IOSConfig.AppKey
	timeInt64 := time.Now().Unix()
	timestamp := strconv.FormatInt(timeInt64, 10)
	messageIOS.Timestamp = timestamp
	// 通过判断是否设置 deviceToken，来区分 单播 和 全播
	if deviceToken == "" {
		// 全播
		messageIOS.Type = "broadcast"
	} else {
		// 单播
		messageIOS.Type = "unicast"
		messageIOS.DeviceTokens = deviceToken
	}
	messageIOS.ProductionMode = Config.PushProductionMode
	messageIOS.Description = contentAvailable
	postBody, _ := json.Marshal(messageIOS)
	url := Config.HostUrl + Config.PostUri
	// MD5加密
	sign := Md5("POST" + url + string(postBody) + Config.IOSConfig.MasterSecret)
	url = url + "?sign=" + sign
	var timeout time.Duration
	res, err = PostJSON(context.Background(), url, messageIOS, map[string]string{}, timeout)
	return
}
