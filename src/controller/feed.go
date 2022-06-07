package controller

import (
	"douyin/src/repository"
	"douyin/src/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	//判断时间戳是否合法
	latestTime := service.TimeStamp(c.Query("latest_time"))
	//判断token是否合法
	tokenUserId := service.Token2Id(c.Query("token"))
	//获取视频列表
	videoList_db := repository.VideoQueryByLastTime(latestTime)
	//返回响应
	c.JSON(http.StatusOK, service.FeedResponse{
		Response:  service.Response{StatusCode: 0},
		VideoList: service.FromDBVideosToMesVideos(videoList_db, tokenUserId),
		NextTime:  time.Now().UnixMilli(),
	})
}
