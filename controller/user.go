package controller

import (
	"douyin/common"
	dtb "douyin/database"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//根据用户id获取用户信息
func UserInfo(c *gin.Context) {

	var token_user_id int64
	tokenString := c.Query("token")
	//查看token是否合法
	if id, isValid := common.TokenValidity(tokenString); !isValid {
		token_user_id = 0
	} else {
		//合法的token解析出登陆用户ID并赋值给token_user_id变量
		token_user_id = id
	}

	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	//根据userid查找user
	if user_db, exist := dtb.UserQueryByID(common.DB, user_id); exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     FromDBUsersTOMesUsers(common.DB, user_db, token_user_id)[0],
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}

}

//用户注册
func Register(c *gin.Context) {

	username := c.Query("username")

	password := Encryption(c.Query("password"))

	if _, exist := dtb.UserQueryByName(common.DB, username); exist {
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
			Avatar:    avator,
			Signature: signature,
		}
		err := dtb.UserCreate(common.DB, &user)
		tokenString, err := GetToken(username, user.Id)
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  fmt.Sprintf("Token release error:%v", err),
				}})
			return
		}
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  fmt.Sprintf("User regist error:%v", err),
				}})
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: 0,
					StatusMsg:  "User regist sucess",
				},
				UserId: user.Id,
				Token:  tokenString,
			})
		}
	}

}
func Login(c *gin.Context) {

	username := c.Query("username")
	password := Encryption(c.Query("password"))
	if user, exist := dtb.UserValid(common.DB, username, password); exist {
		//每次登陆都会更新token
		tokenString, err := GetToken(username, user[0].Id)
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  fmt.Sprintf("Token release error:%v", err),
				}})
			return
		}

		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  fmt.Sprintf("Token update error:%v", err),
				},
			})
			return
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user[0].Id,
			Token:    tokenString,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}

}
