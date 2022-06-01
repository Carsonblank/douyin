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
)

func Publish(c *gin.Context) {

	token := c.PostForm("token")
	var token_user_id int64
	//查看token是否合法
	if id, isValid := common.TokenValidity(token); !isValid {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Token is invalid.",
		},
		)
		//对于视频投稿操作，token非法直接返回
		return
	} else {
		//合法的token解析出登陆用户ID并赋值给token_user_id变量
		token_user_id = id
		//查看合法token返回的用户ID是否存在
		if _, exist := dtb.UserQueryByID(common.DB, token_user_id); !exist {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
			)
			return
		}
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
	finalName := fmt.Sprintf("%d_%s", token_user_id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//更新Video数据库
	if err := dtb.VideoCreate(common.DB, dtb.Video{
		UserId:   token_user_id,
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
func PublishList(c *gin.Context) {

	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	videoList := dtb.VideoQueryByUserID(common.DB, user_id)

	//获取登录用户的id
	token := c.Query("token")
	var token_user_id int64
	//查看token是否合法
	if id, isValid := common.TokenValidity(token); !isValid {
		token_user_id = 0
	} else {
		//合法的token解析出登陆用户ID并赋值给token_user_id变量
		token_user_id = id
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: FromDBVideosToMesVideos(common.DB, videoList, token_user_id),
	})

}
