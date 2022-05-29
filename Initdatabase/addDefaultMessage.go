package main

import (
	"douyin/config"
	dtb "douyin/database"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		log.Panic("Database connect error : ", err)
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&dtb.User{})
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&dtb.Video{})
	if _, exist := dtb.UserQueryByName(db, "DefaultUser"); !exist {
		user := dtb.User{
			Name:      "DefaultUser",
			Password:  "123456",
			Token:     "DefaultUser123456",
			Avatar:    "http://192.168.139.131:8080/static/defaultAvatar.png",
			Signature: "Hello douyin",
		}
		if err := dtb.UserCreate(db, &user); err != nil {
			fmt.Printf("Create user error: %v", err)
		}

		//user, _ := dtb.UserQueryByName(db, "DefaultUser")
		if err := dtb.VideoCreate(db, dtb.Video{
			UserId:   user.Id,
			PlayUrl:  "http://192.168.139.131:8080/static/AttackOnTitan.mp4",
			CoverUrl: "http://192.168.139.131:8080/static/AttackOnTitan.png",
			Title:    "进击的巨人最终季part-2: 2",
		}); err != nil {
			fmt.Printf("Create video error: %v", err)
		}
	}
}
