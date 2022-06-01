package controller

import (
	"douyin/common"
	dtb "douyin/database"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Feed same demo video list for every request
func Feed(c *gin.Context) {

	//因为Query接收的是string，转为int64
	latest_time, _ := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	//判断是否有值
	if latest_time == 0 {
		latest_time = time.Now().UnixMilli()
	}
	videoList_db := dtb.VideoQueryByLastTime(common.DB, latest_time) //last_time)

	//根据token判断是不是登录用户，是登陆用户查找Relation关系表，看是不是关注的用户
	token := c.Query("token")
	var token_user_id int64
	//查看token是否合法
	if id, isValid := common.TokenValidity(token); !isValid {
		token_user_id = 0
	} else {
		//合法的token解析出登陆用户ID并赋值给token_user_id变量
		token_user_id = id
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: FromDBVideosToMesVideos(common.DB, videoList_db, token_user_id),
		NextTime:  time.Now().UnixMilli(),
	})

}
