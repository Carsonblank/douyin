package controller

import (
	dtb "douyin/database"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//根据用户id获取用户信息
func UserInfo(db *gorm.DB) gin.HandlerFunc {
	fun := func(c *gin.Context) {
		var token_user_id int64
		if users, exit := dtb.UserQueryByToken(db, c.Query("token")); exit {
			token_user_id = users[0].Id
		} else {
			token_user_id = 0
		}

		user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
		//根据userid查找user
		if user_db, exist := dtb.UserQueryByID(db, user_id); exist {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 0},
				User:     FromDBUsersTOMesUsers(db, user_db, token_user_id)[0],
			})
		} else {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
		}
	}
	return gin.HandlerFunc(fun)
}

//用户注册
func Register(db *gorm.DB) gin.HandlerFunc {
	fun := func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("password")

		if _, exist := dtb.UserQueryByName(db, username); exist {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
			})
		} else {
			avator := c.Query("avator")
			signature := c.Query("signature")
			if avator == "" {
				avator = "http://192.168.139.131:8080/static/defaultAvatar.png"
			}
			if signature == "" {
				signature = "Hello douyin"
			}
			user := dtb.User{
				Name:      username,
				Password:  password,
				Token:     GetToken(username, password),
				Avatar:    avator,
				Signature: signature,
			}
			err := dtb.UserCreate(db, &user)
			if err != nil {
				log.Printf("User create error : %v \n", err)
			} else {
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{
						StatusCode: 0,
						StatusMsg:  "User regist sucess",
					},
					UserId: user.Id,
					Token:  user.Token,
				})
			}
		}
	}
	return gin.HandlerFunc(fun)
}
func Login(db *gorm.DB) gin.HandlerFunc {
	fun := func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("password")
		if user, exist := dtb.UserValid(db, username, password); exist {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   user[0].Id,
				Token:    user[0].Token,
			})
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
		}

	}
	return gin.HandlerFunc(fun)
}
