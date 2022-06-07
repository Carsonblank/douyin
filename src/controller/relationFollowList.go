package controller

import (
	"douyin/src/repository"
	"douyin/src/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FollowerList all users have same follow list
func FollowerList(c *gin.Context) {
	//验证token合法性
	tokenUserId := service.Token2Id(c.Query("token"))
	//验证userId合法性
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if _, exist := repository.UserQueryByID(userId); !exist {
		c.JSON(http.StatusOK, service.UserListResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  "User is not exist!\n",
			},
		})
		return
	}
	//获取粉丝列表
	userList_db := repository.UserQueryByFollowerID(userId)
	//返回成功响应
	c.JSON(http.StatusOK, service.UserListResponse{
		Response: service.Response{
			StatusCode: 0,
		},
		UserList: service.FromDBUsersTOMesUsers(userList_db, tokenUserId),
	})

}
