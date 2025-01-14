package install

import (
	"first/init/mysql"
	"first/internal/model"
)

// DBInstall 数据迁移
func DBInstall() {
	var err error
	db := mysql.NewDB()

	// err = db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model.ModelList...)
	err = db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate()
	if err != nil {
		panic("table migration err")
	}
	modelList := model.ModelList
	for _, v := range modelList {
		thisModel := v.(model.MonitorModel)
		db.Exec("ALTER TABLE `" + thisModel.TableName() + "` COMMENT '" + thisModel.Comment() + "'")
	}
	// 设备数据表
	sql := `
CREATE TABLE 
IF 
NOT EXISTS device_data ( 
	id bigint(20) NOT NULL, 
	device_id BIGINT ( 20 ) NOT NULL DEFAULT 0 COMMENT '设备id',
	name VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT 'eg wind_speed',
	port VARCHAR ( 36 ) NOT NULL DEFAULT '' COMMENT '接口', 
	value double(20,2) NOT NULL DEFAULT 0 COMMENT '值',
	original_value double(20,2) NOT NULL DEFAULT -99 COMMENT '值',
	created_time BIGINT ( 20 ) UNSIGNED NOT NULL DEFAULT 0 COMMENT '时间戳', 
	PRIMARY KEY ( id, device_id ), 
	KEY idx_d_n_c ( device_id,name,created_time),
	KEY idx_d_c ( device_id,created_time)
) ENGINE = INNODB COMMENT = '设备数据表' PARTITION BY HASH ( device_id ) PARTITIONS 1024 
`
	if err := db.Exec(sql).Error; err != nil {
		panic("table migration device_data err")
	}
}
