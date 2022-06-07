package controller

import (
	"douyin/src/database"
	"douyin/src/repository"
	"douyin/src/service"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func Publish(c *gin.Context) {
	//解析token，查看token是否合法
	userId := service.Token2Id(c.PostForm("token"))
	if userId == 0 {
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  "Token is invalid.",
		},
		)
		return
	}
	//查看解析出来的用户是否存在
	if _, exist := repository.UserQueryByID(userId); !exist {
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		},
		)
		return
	}
	//接收文件
	data, err := c.FormFile("data")
	//文件接受错误
	if err != nil {
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//获取传输文件名
	filename := filepath.Base(data.Filename)
	//由用户ID+时间戳+传输文件名，构造最终的文件名
	finalName := fmt.Sprintf("%d_%d_%s", userId, time.Now().UnixMicro(), filename)
	//添加保存路径
	saveFile := filepath.Join("./public/", finalName)
	//保存文件
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		//文件存储错误，返回错误响应
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//使用ffmpeg获取封面
	coverName := service.FfmpegCreateCover(saveFile)
	//更新Video数据库
	if err := repository.VideoCreate(database.Video{
		UserId:   userId,
		PlayUrl:  "http://192.168.139.131:8080/static/" + finalName,
		CoverUrl: "http://192.168.139.131:8080/static/" + filepath.Base(coverName),
		Title:    c.PostForm("title"),
	}); err != nil {
		//数据库对象创建失败，返回错误响应
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	} else {
		//创建成功，返回成功响应
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 0,
			StatusMsg:  "Video publish success!",
		})
	}

}
