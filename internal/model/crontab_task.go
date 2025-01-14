package model

type CrontabTask struct {
	ID               string `gorm:"column:id;type:varchar(30);not null;comment:主键ID;" json:"id"`
	Name             string `gorm:"column:name;type:varchar(50);not null;default:'';comment:地区名称[一般同area];" json:"name"`
	Func             string `gorm:"column:func;type:varchar(100);not null;default:'';comment:调用接口名称;" json:"func"`
	Cron             string `gorm:"column:cron;type:varchar(100);not null;default:'';comment:cron执行表达式;" json:"cron"`
	Status           int8   `gorm:"column:status;type:tinyint(2);not null;default:1;comment:启用状态 -1已删除 1.启用 2:禁用;" json:"status"`
	Remark           string `gorm:"column:remark;type:varchar(100);not null;default:'';comment:备注信息;" json:"remark"`
	LastIssuTime     int64  `gorm:"column:last_issu_time;type:bigint(20);not null;default:0;comment:上次执行时间;" json:"last_issu_time"`
	LastIssuResponse string `gorm:"column:last_issu_response;type:text;not null;comment:上次执行响应;" json:"last_issu_response"`
	CreatedTime      int64  `gorm:"column:created_time;type:bigint(20);not null;default:0;comment:创建时间;" json:"created_time"`
	UpdatedTime      int64  `gorm:"column:updated_time;type:bigint(20);not null;default:0;comment:更新时间;" json:"updated_time"`
	DeletedTime      int64  `gorm:"column:deleted_time;type:bigint(20);not null;default:0;comment:删除时间;" json:"deleted_time"`
}

func init() {
	ModelList = append(ModelList, &CrontabTask{})
}

// TableName 表名
func (CrontabTask) TableName() string {
	return "crontab_task"
}
