package controller

import (
	"douyin/src/repository"
	"douyin/src/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//获取点赞列表
func FavoriteList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	var exist bool
	//验证Id是否合法
	if _, exist = repository.UserQueryByID(userId); !exist {
		c.JSON(http.StatusOK, service.VideoListResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  "user is not exist.",
			},
		})
		return
	}
	//获取点赞列表视频
	videoList, _ := repository.FavoriteQuerybyUserID(userId)
	//返回成功响应
	c.JSON(http.StatusOK, service.VideoListResponse{
		Response: service.Response{
			StatusCode: 0,
		},
		VideoList: service.FromDBVideosToMesVideos(videoList, userId),
	})

}
