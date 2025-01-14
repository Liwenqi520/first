package mysql

import (
	"fmt"

	"github.com/Liwenqi520/gormx"
	"github.com/Liwenqi520/logger"
	"gorm.io/gorm"
)

var _db *gorm.DB

// MysqlInit mysql数据库初始化
// 失败重连5次
func MysqlInit(cfg gormx.Config, logCfg logger.Config) (err error) {
	_db, err = gormx.New(cfg)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// 获取一个数据库对象
func NewDB() *gorm.DB {
	return _db
}
