package main

import (
	"demo/src/config"
	"demo/src/controller"
	"github.com/gin-gonic/gin"
	"log"
)

func initRouter(r *gin.Engine) {

	apiRouter := r.Group("/douyin")

	// basic apis

	apiRouter.POST("/user/register/", controller.Register)

}

func Run(httpServer *gin.Engine) {
	initRouter(httpServer)
	//log.Fatal会在里面函数运行出现错误时打印日志
	log.Fatal(httpServer.Run(config.GetServerConfig().HTTP_HOST + ":" + config.GetServerConfig().HTTP_PORT))
}
