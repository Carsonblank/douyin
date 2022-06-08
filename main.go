package main

import (
	"douyin/config"
	"douyin/src/controller"
	"douyin/src/database"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	//获取数据库配置信息
	config.GetMySQLConfig()
	//服务器配置
	config.GetServerConfig()
	//其他信息
	config.GetFeedNums()
}

func main() {

	database.InitDatabase()
	//设置release模式
	gin.SetMode(config.GIN_MODE)
	//设置debug模式
	//gin.SetMode(gin.DebugMode)
	//开启高亮
	gin.ForceConsoleColor()
	//关闭高亮
	//gin.DisableConsoleColor()
	//输出日志
	f, _ := os.Create(fmt.Sprintf("./log/gin_%v_%v_%v%s", time.Now().Year(), time.Now().Month(), time.Now().Day(), ".log"))
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()

	initRouter(r)

	r.Run(config.DouyinPort) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)
	//
	//// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)
	//
	//// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)

}
