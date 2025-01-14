package model

//Transcation 用户交易流水表
type Transcation struct {
	ID           string `gorm:"column:id;type:varchar(30);NOT NULL;comment:主键;" json:"id"`
	UserID       string `gorm:"column:user_id;type:varchar(30);not null;default:'';comment:用户id;" json:"user_id"`
	TotalAmount  int64  `gorm:"column:total_amount;type:int;not null;default:0;comment:用户当前总金额;" json:"total_amount"`
	ChangeAmount int64  `gorm:"column:change_amount;type:int;not null;default:0;comment:用户已变动的金额;" json:"change_amount"`
	CreatedTime  int64  `gorm:"column:created_time;type:bigint(20);not null;default:0;comment:创建时间;" json:"created_time"`
	UpdatedTime  int64  `gorm:"column:updated_time;type:bigint(20);not null;default:0;comment:更新时间;" json:"updated_time"`
	DeletedTime  int64  `gorm:"column:deleted_time;type:bigint(20);not null;default:0;comment:删除时间;" json:"deleted_time"`
}

func init() {
	ModelList = append(ModelList, &Transcation{})
}

// TableName 表名
func (Transcation) TableName() string {
	return "transcation"
}
