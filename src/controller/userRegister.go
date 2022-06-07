package controller

import (
	"douyin/src/database"
	"douyin/src/repository"
	"douyin/src/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//用户注册
func Register(c *gin.Context) {

	username := c.Query("username")
	//注册用户是否存在
	if _, exist := repository.UserQueryByName(username); exist {
		//返回注册失败响应
		c.JSON(http.StatusOK, service.UserLoginResponse{
			Response: service.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})

		return
	}
	//散列加密密码
	password := service.Encryption(c.Query("password"))
	//是否有给定头像和签名
	avator := c.Query("avator")
	signature := c.Query("signature")
	if avator == "" {
		avator = "http://192.168.139.131:8080/static/defaultAvatar.png"
	}
	if signature == "" {
		signature = "Hello douyin"
	}

	//创建用户
	user := database.User{
		Name:      username,
		Password:  password,
		Avatar:    avator,
		Signature: signature,
	}
	err := repository.UserCreate(&user)
	//用户创建失败
	if err != nil {
		//返回注册失败响应
		c.JSON(http.StatusOK, service.UserLoginResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("User regist error:%v", err),
			}})
		return
	}

	//根据用户名和Id创建token
	tokenString, err := service.GetToken(username, user.Id)
	//token创建失败
	if err != nil {
		//返回token生成失败响应
		c.JSON(http.StatusOK, service.UserLoginResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("Token release error:%v", err),
			}})
		return
	}

	//返回用户注册成功响应
	c.JSON(http.StatusOK, service.UserLoginResponse{
		Response: service.Response{
			StatusCode: 0,
			StatusMsg:  "User regist sucess",
		},
		UserId: user.Id,
		Token:  tokenString,
	})

}
