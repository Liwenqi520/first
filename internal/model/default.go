package model

import (
	"first/internal/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MonitorModel interface {
	Comment() string
	TableName() string
}

//ModelList 模型列表
var ModelList []interface{}

//ID未删除
func WithID(ID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ? and deleted_time = 0", ID)
	}
}

//ID全部
func WithIDTotal(ID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", ID)
	}
}

func WithCompany(userSession common.UserInfo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("company_id = ? and deleted_time = 0", userSession.CompanyID)
	}
}

func WithIDAndCompany(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	userSession, _ := common.SessionGet(c)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", c.Param("id")).
			Where("company_id = ? and deleted_time = 0", userSession.CompanyID)
	}
}

//分页
// func Paginate(PageInfo *schema.PageResult) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		return db.Clauses(clause.Limit{
// 			Limit:  PageInfo.Limit,
// 			Offset: PageInfo.Offset,
// 		})
// 	}
// }
