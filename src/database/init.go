package database

import (
	"douyin/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//全局变量
var MySqlDb *gorm.DB

//初始化数据库
func InitDatabase() {

	//打开数据库
	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		log.Panic("Database connect error : ", err)
	}

	//自动迁移表
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(User{})
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(Video{})
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(Favorite{})
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(Comment{})
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(Relation{})

	//赋值给全局变量
	MySqlDb = db
}
