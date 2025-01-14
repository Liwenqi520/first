package migrate

import (
	"first/init/mysql"
	"first/internal/model"
	"fmt"
)

type versionScript struct {
	DstVersion string       //要升级的最新版
	Func       func() error //升级脚本
}

var versions = map[string]versionScript{
	"1": {
		DstVersion: "2",
		Func:       Migrate2,
	},
	// "2": {
	// 	DstVersion: "3",
	// 	Func:       Migrate3,
	// },
}

func Migrate() {
	var (
		currentVer model.Version
		newVer     model.Version
	)
	db := mysql.NewDB()
	db.Model(&model.Version{}).Take(&currentVer)
	var (
		v = currentVer.Version
	)
	for {
		script, ok := versions[v]
		if ok && script.DstVersion != "" {
			err := script.Func()
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			newVer.Version = script.DstVersion
			db.Model(&model.Version{}).Where("version = ?", v).Updates(&newVer)
			v = script.DstVersion
		} else {
			break
		}
	}
}
