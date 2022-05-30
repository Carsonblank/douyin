//关注操作
package controller

import (
	"douyin/common"
	dtb "douyin/database"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RelationAction no practical effect, just check if token is valid
func RelationAction(db *gorm.DB) gin.HandlerFunc {
	fun := func(c *gin.Context) {
		token := c.Query("token")

		//查看token是否合法
		if !common.TokenValidity(token) {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Token is invalid.",
			},
			)
			return
		}

		if user, exist := dtb.UserQueryByToken(db, token); exist {
			to_user_id, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
			if _, exist := dtb.UserQueryByID(db, to_user_id); !exist {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "To user doesn't exist"})
				return
			}
			//
			action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
			if action_type == 1 {
				//防止重复关注
				if dtb.FavoriteQueryByUserAndUser(db, user[0].Id, to_user_id) {
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Repeat follow."})
				} else {
					if err := dtb.RelationCreate(db, dtb.Relation{
						UserId:   user[0].Id,
						ToUserId: to_user_id,
					}); err != nil {
						c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: fmt.Sprintf("Relation create error: %v", err)})
					} else {
						c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Relation success"})
					}
				}

			} else if action_type == 2 {
				if !dtb.FavoriteQueryByUserAndUser(db, user[0].Id, to_user_id) {
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Repeat delete follow."})
					return
				}

				if err := dtb.RelationDelete(db, user[0].Id, to_user_id); err != nil {
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: fmt.Sprintf("Relation delete error: %v", err)})
				} else {
					c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Delete relation success\n"})
				}

			} else {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Message action_type must be 1 or 2"})
			}

		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		}

	}
	return gin.HandlerFunc(fun)
}

// FollowList all users have same follow list
func FollowList(db *gorm.DB) gin.HandlerFunc {
	fun := func(c *gin.Context) {
		//
		token := c.Query("token")
		var token_user_id int64
		if user, exist := dtb.UserQueryByToken(db, token); exist {
			token_user_id = user[0].Id
		} else {
			token_user_id = 0
		}

		//
		user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
		if _, exist := dtb.UserQueryByID(db, user_id); !exist {
			c.JSON(http.StatusOK, UserListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "User is not exist!\n",
				},
			})
			return
		}
		userList_db := dtb.UserQueryByFollowID(db, user_id)

		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: FromDBUsersTOMesUsers(db, userList_db, token_user_id),
		})
	}
	return gin.HandlerFunc(fun)
}

// FollowerList all users have same follow list
func FollowerList(db *gorm.DB) gin.HandlerFunc {
	fun := func(c *gin.Context) {
		token := c.Query("token")
		var token_user_id int64
		if user, exist := dtb.UserQueryByToken(db, token); exist {
			token_user_id = user[0].Id
		} else {
			token_user_id = 0
		}

		//
		user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
		if _, exist := dtb.UserQueryByID(db, user_id); !exist {
			c.JSON(http.StatusOK, UserListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "User is not exist!\n",
				},
			})
			return
		}
		userList_db := dtb.UserQueryByFollowerID(db, user_id)

		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserList: FromDBUsersTOMesUsers(db, userList_db, token_user_id),
		})
	}
	return gin.HandlerFunc(fun)
}
