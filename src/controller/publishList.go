package controller

import (
	"douyin/src/repository"
	"douyin/src/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PublishList(c *gin.Context) {

	//获取将要查看用户视频列表的用户Id
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	videoList := repository.VideoQueryByUserID(userId)

	//获取登录用户的Id
	tokenUserId := service.Token2Id(c.Query("token"))

	//返回视频列表
	c.JSON(http.StatusOK, service.VideoListResponse{
		Response: service.Response{
			StatusCode: 0,
		},
		VideoList: service.FromDBVideosToMesVideos(videoList, tokenUserId),
	})

}
