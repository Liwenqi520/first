package model

type Version struct {
	Version string `json:"version" gorm:"column:version;type:varchar(50) NOT NULL;comment:版本号"`
}

func init() {
	ModelList = append(ModelList, &Version{})
}

func (Version) Comment() string {
	return "版本表"
}

func (Version) TableName() string {
	return "version"
}
