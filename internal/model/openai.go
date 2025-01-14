package model

//OpenAILog OpenAILog日志表
type OpenaiLog struct {
	ID          int64  `gorm:"column:id;type:bigint;AUTO_INCREMENT;primary_key;comment:主键;" json:"id"`
	UserID      string `gorm:"column:user_id;type:varchar(30);not null;default:'';comment:用户id;" json:"user_id"`
	ReqParams   string `gorm:"column:req_params;type:text;comment:请求参数;" json:"req_params"`
	Response    string `gorm:"column:response;type:text;comment:响应参数;" json:"response"`
	WasteTime   int64  `gorm:"column:waste_time;type:int;default:0;comment:耗时;" json:"waste_time"`
	UsedToken   int64  `gorm:"column:used_token;type:int;default:0;comment:使用token数;" json:"used_token"`
	CreatedTime int64  `gorm:"column:created_time;type:bigint(20);default:0;comment:创建时间;" json:"created_time"`
	UpdatedTime int64  `gorm:"column:updated_time;type:bigint(20);default:0;comment:更新时间;" json:"updated_time"`
	DeletedTime int64  `gorm:"column:deleted_time;type:bigint(20);default:0;comment:删除时间;" json:"deleted_time"`
}

func init() {
	ModelList = append(ModelList, &OpenaiLog{})
}

// TableName 表名
func (OpenaiLog) TableName() string {
	return "openai_log"
}
