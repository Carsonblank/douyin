package middleware

import (
	"github.com/RaymondCode/simple-demo/src/common"
	"github.com/RaymondCode/simple-demo/src/controller"
	"github.com/RaymondCode/simple-demo/src/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取参数
		userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, controller.UserLoginResponse{
				Response: controller.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
			c.Abort()
			return
		}
		tokenString := c.Query("token")

		//验证参数
		if tokenString == "" {
			c.JSON(http.StatusOK, controller.UserLoginResponse{
				Response: controller.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
			c.Abort()
			return
		}
		//验证token
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusOK, controller.UserLoginResponse{
				Response: controller.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
			c.Abort()
			return
		}

		//验证通过后获取claim中的userId
		username := claims.Username
		userInfo, err := service.QueryUserInfo(username)
		if err != nil || userInfo == nil {
			log.Println(err)
			c.JSON(http.StatusOK, controller.UserLoginResponse{
				Response: controller.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
			c.Abort()
			return
		}
		if userInfo.User.Id != userId {
			log.Println(err)
			c.JSON(http.StatusOK, controller.UserLoginResponse{
				Response: controller.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
			c.Abort()
			return
		}
		c.Set("userInfo", userInfo)
		c.Next()
	}

}
