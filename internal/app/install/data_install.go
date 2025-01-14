package install

import (
	"first/init/mysql"
	"first/internal/defines"
	"first/internal/model"
)

func DataInstall() {

	// 插入版本信息
	var ver model.Version
	mysql.NewDB().Model(&model.Version{}).Attrs(model.Version{
		Version: defines.DBVersion,
	}).FirstOrCreate(&ver)
}
