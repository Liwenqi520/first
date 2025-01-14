package model

// SendMsgLog 消息发送记录
type SendMsgLog struct {
	ID          string `gorm:"column:id;type:varchar(20);comment:主键;primaryKey;" json:"id"`
	CompanyID   string `gorm:"column:company_id;type:varchar(30);default:'';not null;comment:企业ID" json:"company_id"`
	SendType    int8   `gorm:"column:send_type;type:tinyint(2);not null;default:1;comment:发送类型1邮箱 2短信;" json:"send_type"`
	Target      string `gorm:"column:target;type:varchar(255);not null;default:'';comment:介质" json:"target"`
	Content     string `gorm:"column:content;type:varchar(2550);not null;default:'';comment:发送内容;" json:"content"`
	Status      int8   `gorm:"column:status;type:tinyint(2);not null;default:0;comment:发送状态0待发送 1成功 2失败;" json:"status"`
	Res         string `gorm:"column:res;type:varchar(2550);not null;default:'';comment:发送返回值;" json:"res"`
	CreatedTime int64  `gorm:"column:created_time;not null;default:0;comment:创建的时间戳" json:"created_time"`
	UpdatedTime int64  `gorm:"column:updated_time;not null;default:0;comment:更新的时间戳" json:"updated_time"`
	DeletedTime int64  `gorm:"column:deleted_time;not null;default:0;comment:删除的时间戳" json:"deleted_time"`
}

func init() {
	ModelList = append(ModelList, &SendMsgLog{})
}

func (SendMsgLog) Comment() string {
	return "消息发送记录"
}

func (SendMsgLog) TableName() string {
	return "send_msg_log"
}
