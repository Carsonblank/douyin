package controller

import (
	"douyin/src/repository"
	"douyin/src/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//根据用户id获取用户信息
func UserInfo(c *gin.Context) {

	//查看token是否合法
	token_user_id := service.Token2Id(c.Query("token"))

	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	//根据userid查找user
	if user_db, exist := repository.UserQueryByID(user_id); exist {
		//找到用户信息，返回查找成功响应
		c.JSON(http.StatusOK, service.UserResponse{
			Response: service.Response{StatusCode: 0},
			User:     service.FromDBUsersTOMesUsers(user_db, token_user_id)[0],
		})
	} else {
		//未找到用户信息，返回查找失败响应
		c.JSON(http.StatusOK, service.UserResponse{
			Response: service.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}

}
