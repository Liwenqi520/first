package service

import (
	"first/init/mysql"
	"first/internal/model"
	"fmt"
)

// 初始化所有定时任务入口
func CronInit() (err error) {
	fmt.Println("cronInit start~")
	// 取出数据表中的定时任务
	db := mysql.NewDB()
	var CrontabTaskList []model.CrontabTask
	db.Model(&model.CrontabTask{}).Where("status = 1 and deleted_time = 0").Find(&CrontabTaskList)
	if len(CrontabTaskList) > 0 {
		for _, task := range CrontabTaskList {
			CronLog{}.StartCronCommon(task.Cron, task.Func)
		}
	}
	return nil
}
