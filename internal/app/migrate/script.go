package migrate

import (
	"first/init/mysql"
	v2 "first/internal/model/history/v2"
)

func Migrate2() error {
	db := mysql.NewDB()
	// eg:
	db.Migrator().AddColumn(&v2.AuthUser{}, "remark")
	db.Migrator().AlterColumn(&v2.AuthUser{}, "company_pic")
	AuthUserExist := db.Migrator().HasTable(v2.AuthUser{}.TableName())
	if !AuthUserExist {
		db.Migrator().AutoMigrate(&v2.AuthUser{})
	}
	UserCashExist := db.Migrator().HasTable(v2.UserCash{}.TableName())
	if !UserCashExist {
		db.Migrator().AutoMigrate(&v2.UserCash{})
	}
	TranscationExist := db.Migrator().HasTable(v2.Transcation{}.TableName())
	if !TranscationExist {
		db.Migrator().AutoMigrate(&v2.Transcation{})
	}
	OpenaiLogExist := db.Migrator().HasTable(v2.OpenaiLog{}.TableName())
	if !OpenaiLogExist {
		db.Migrator().AutoMigrate(&v2.OpenaiLog{})
	}
	return nil
}
