package controller

import (
	"demo/src/common"
	"demo/src/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//调用service层注册服务，返回id
	userId, err := service.CreateUser(username, password)
	if err != nil {
		log.Println("insert err", err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	}
	//数据库验证通过，发放token
	token, err := common.ReleaseToken(username)
	if err != nil || token == "" {
		log.Println(err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   userId,
		Token:    token,
	})

}
