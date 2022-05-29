package main

import (
	"log"

	"douyin/config"
	"douyin/controller"
	dtb "douyin/database" //database model type

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	db := initDatabase()

	r := gin.Default()

	initRouter(r, db)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

//初始化数据库
func initDatabase() (db *gorm.DB) {

	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		log.Panic("Database connect error : ", err)
	}
	db.Migrator().CreateTable()

	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&dtb.User{})
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&dtb.Video{})
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&dtb.Favorite{})
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&dtb.Comment{})
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&dtb.Relation{})
	return
}

func initRouter(r *gin.Engine, db *gorm.DB) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed(db))
	apiRouter.GET("/user/", controller.UserInfo(db))
	apiRouter.POST("/user/register/", controller.Register(db))
	apiRouter.POST("/user/login/", controller.Login(db))
	apiRouter.POST("/publish/action/", controller.Publish(db))
	apiRouter.GET("/publish/list/", controller.PublishList(db))
	//
	//// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction(db))
	apiRouter.GET("/favorite/list/", controller.FavoriteList(db))
	apiRouter.POST("/comment/action/", controller.CommentAction(db))
	apiRouter.GET("/comment/list/", controller.CommentList(db))
	//
	//// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction(db))
	apiRouter.GET("/relation/follow/list/", controller.FollowList(db))
	apiRouter.GET("/relation/follower/list/", controller.FollowerList(db))

}
