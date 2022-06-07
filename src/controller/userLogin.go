package controller

import (
	"douyin/src/database"
	"douyin/src/repository"
	"douyin/src/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	username := c.Query("username")
	password := service.Encryption(c.Query("password"))
	var user []database.User
	var exist bool
	if user, exist = repository.UserValid(username, password); !exist {
		//如果用户不存在，返回登录失败响应
		c.JSON(http.StatusOK, service.UserLoginResponse{
			Response: service.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}

	//获取新的token，每次登陆都会更新token
	tokenString, err := service.GetToken(username, user[0].Id)
	//token生成失败
	if err != nil {
		//返回登录失败响应
		c.JSON(http.StatusOK, service.UserLoginResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("Token release error:%v", err),
			}})
		return
	}

	//返回登录成功响应
	c.JSON(http.StatusOK, service.UserLoginResponse{
		Response: service.Response{StatusCode: 0},
		UserId:   user[0].Id,
		Token:    tokenString,
	})

}
