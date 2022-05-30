package controller

import (
	"douyin/common"
	dtb "douyin/database"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Publish(db *gorm.DB) gin.HandlerFunc {
	fun := func(c *gin.Context) {
		token := c.PostForm("token")

		//查看token是否合法
		if !common.TokenValidity(token) {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Token is invalid.",
			},
			)
			return
		}
		if _, exist := dtb.UserQueryByToken(db, token); !exist {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
			return
		}

		data, err := c.FormFile("data")
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}

		filename := filepath.Base(data.Filename)
		user, _ := dtb.UserQueryByToken(db, token)
		finalName := fmt.Sprintf("%d_%s", user[0].Id, filename)
		saveFile := filepath.Join("./public/", finalName)
		if err := c.SaveUploadedFile(data, saveFile); err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		//更新Video数据库
		if err := dtb.VideoCreate(db, dtb.Video{
			UserId:   user[0].Id,
			PlayUrl:  "http://192.168.139.131:8080/static/" + finalName,
			CoverUrl: "http://192.168.139.131:8080/static/defaultVideoCover.png", //使用默认封面
			Title:    c.PostForm("title"),
		}); err != nil {
			log.Printf("Video create error: %v\n", err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, Response{
				StatusCode: 0,
				StatusMsg:  "Video publish success!",
			})
		}

	}
	return gin.HandlerFunc(fun)
}
func PublishList(db *gorm.DB) gin.HandlerFunc {
	fun := func(c *gin.Context) {

		user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
		videoList := dtb.VideoQueryByUserID(db, user_id)

		//获取登录用户的id
		token := c.Query("token")
		if token == "" {
			user_id = 0
		} else {
			if user, exist := dtb.UserQueryByToken(db, token); exist {
				user_id = user[0].Id
			} else {
				user_id = 0
			}

		}

		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: FromDBVideosToMesVideos(db, videoList, user_id),
		})
	}
	return gin.HandlerFunc(fun)
}
