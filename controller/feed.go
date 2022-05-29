package controller

import (
	dtb "douyin/database"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Feed same demo video list for every request
func Feed(db *gorm.DB) gin.HandlerFunc {
	fun := func(c *gin.Context) {
		//因为Query接收的是string，转为int64
		latest_time, _ := strconv.ParseInt(c.Query("latest_time"), 10, 64)
		//判断是否有值
		if latest_time == 0 {
			latest_time = time.Now().UnixMilli()
		}
		videoList_db := dtb.VideoQueryByLastTime(db, latest_time) //last_time)

		//根据token判断是不是登录用户，是登陆用户查找Relation关系表，看是不是关注的用户
		token := c.Query("token")
		var user_id int64
		if token == "" {
			user_id = 0
		} else {
			if user, exist := dtb.UserQueryByToken(db, token); exist {
				user_id = user[0].Id
			} else {
				user_id = 0
			}
		}
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: FromDBVideosToMesVideos(db, videoList_db, user_id),
			NextTime:  time.Now().UnixMilli(),
		})

	}
	return gin.HandlerFunc(fun)
}
