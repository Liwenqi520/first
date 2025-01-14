package model

//UserCash 用户余额表
type UserCash struct {
	ID              string `gorm:"column:id;type:varchar(30);NOT NULL;comment:主键;" json:"id"`
	UserID          string `gorm:"column:user_id;type:varchar(30);not null;default:'';comment:用户id;" json:"user_id"`
	TotalAmount     int64  `gorm:"column:total_amount;type:int;not null;default:0;comment:用户当前总金额;" json:"total_amount"`
	UsedAmount      int64  `gorm:"column:used_amount;type:int;not null;default:0;comment:用户已使用金额;" json:"used_amount"`
	AvailableAmount int64  `gorm:"column:available_amount;type:int;not null;default:0;comment:用户可用金额;" json:"available_amount"`
	FrozenAmount    int64  `gorm:"column:frozen_amount;type:int;not null;default:0;comment:用户冻结金额;" json:"frozen_amount"`
	LastUsedTime    int64  `gorm:"column:last_used_time;type:int;not null;default:0;comment:用户上次使用时间;" json:"last_used_time"`
	CreatedTime     int64  `gorm:"column:created_time;type:bigint(20);not null;default:0;comment:创建时间;" json:"created_time"`
	UpdatedTime     int64  `gorm:"column:updated_time;type:bigint(20);not null;default:0;comment:更新时间;" json:"updated_time"`
	DeletedTime     int64  `gorm:"column:deleted_time;type:bigint(20);not null;default:0;comment:删除时间;" json:"deleted_time"`
}

// TableName 表名
func (UserCash) TableName() string {
	return "user_cash"
}
