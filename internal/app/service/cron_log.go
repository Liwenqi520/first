package service

import (
	"first/init/mysql"
	"first/internal/app/logic"
	"first/internal/model"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
)

type CronLog struct{}

// StartCronCommon 启动定时任务
func (cl CronLog) StartCronCommon(spec, funcName string) error {
	c := cron.New(cron.WithSeconds())
	// 秒 分 时 日 月 周
	// spec := "0 0 8 * * *" //
	ToRun := cl.GetTaskFunc(funcName)
	_, err := c.AddFunc(spec, ToRun)
	if err != nil {
		return err
	}
	c.Start()
	return nil
}

func (cl CronLog) GetTaskFunc(FuncName string) (ToRun func()) {
	switch FuncName {
	case "PushDeviceData":
		ToRun = func() {
			cl.PushDeviceData()
		}

	}
	return
}

// PushDeviceData 推送设备数据
func (cl CronLog) PushDeviceData() {
	// Domain := AppInit.AppConfig.TestPushDevice.Domain
	Config := logic.Common{}.GetPushTestDeviceDataConfig()
	url := Config.Domain + "/device/timing-push"
	res, err := http.Get(url)
	body, err := ioutil.ReadAll(res.Body)
	db := mysql.NewDB()
	var crontabTask model.CrontabTask
	QueryTask := db.Model(&model.CrontabTask{}).Where("func = ? and deleted_time = 0", "PushDeviceData")
	QueryTask.First(&crontabTask)
	UpdateParams := make(map[string]interface{})
	UpdateParams["last_issu_time"] = time.Now().Unix()
	UpdateParams["last_issu_response"] = string(body)
	if err != nil {
		UpdateParams["last_issu_response"] = err.Error()
	}
	db.Model(&model.CrontabTask{}).Where("id = ?", crontabTask.ID).UpdateColumns(&UpdateParams)
}

func (CronLog) LogTimeToUnix(strTime string) int64 {
	if len(strTime) < 20 {
		return 0
	}
	const layout = "2006-01-02T15:04:05"
	t, err := time.ParseInLocation(layout, strTime[0:19], time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

type LogInfo struct {
	Level       string `json:"level"`
	Time        string `json:"time"`
	TimeUnix    int64  `json:"time_unix"`
	Caller      string `json:"caller"`
	Msg         string `json:"msg"`
	Version     string `json:"version"`
	ServerName  string `json:"server_name"`
	TraceID     string `json:"trace_id"`
	LogType     string `json:"log_type"`
	CompanyID   string `json:"company_id"`
	UserID      string `json:"user_id"`
	Tag         string `json:"tag"`          // 标志 tag = user 用户日志
	Data        string `json:"data"`         // 用户操作数据json
	OperateType string `json:"operate_type"` //操作类型 insert,delete,update,current
}
