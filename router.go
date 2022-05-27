package main

import (
	"github.com/RaymondCode/simple-demo/src/config"
	"github.com/RaymondCode/simple-demo/src/controller"
	"github.com/RaymondCode/simple-demo/src/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"path"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")
	r.GET("root/test/simple-demo-main/public/:name", func(c *gin.Context) {
		name := c.Param("name")
		//拼接路径,如果没有这一步，则默认在当前路径下寻找
		filename := path.Join("./public/", name)
		//响应一个文件
		c.File(filename)
		return
	})
	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", middleware.AuthMiddleware(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}

func Run(httpServer *gin.Engine) {
	initRouter(httpServer)
	log.Fatal(httpServer.Run(config.GetServerConfig().HTTP_HOST + ":" + config.GetServerConfig().HTTP_PORT))
}
