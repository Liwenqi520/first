package model

//SystemConfig 系统配置
type SystemConfig struct {
	ID          string `gorm:"column:id;type:varchar(30);NOT NULL;comment:主键;" json:"id"`
	Name        string `gorm:"column:name;type:varchar(100);not null;default:'';comment:配置名称;" json:"name"`
	Key         string `gorm:"column:key;type:varchar(100);not null;default:'';comment:key;" json:"key"`
	Value       string `gorm:"column:value;type:varchar(300);not null;default:'';comment:value;" json:"value"`
	Cate        int64  `gorm:"column:cate;type:tinyint(1);not null;default:0;comment:系统内置 0否 1是;" json:"cate"`
	Remark      string `gorm:"column:remark;type:varchar(200);not null;default:'';comment:备注;" json:"remark"`
	CreatedID   string `gorm:"column:created_id;type:varchar(30);not null;default:'';comment:创建人ID;" json:"created_id"`
	CreatedTime int64  `gorm:"column:created_time;type:bigint(20);not null;default:0;comment:创建时间;" json:"created_time"`
	UpdatedTime int64  `gorm:"column:updated_time;type:bigint(20);not null;default:0;comment:更新时间;" json:"updated_time"`
	DeletedTime int64  `gorm:"column:deleted_time;type:bigint(20);not null;default:0;comment:删除时间;" json:"deleted_time"`
}

func init() {
	ModelList = append(ModelList, &SystemConfig{})
}

// TableName 表名
func (SystemConfig) TableName() string {
	return "system_config"
}
