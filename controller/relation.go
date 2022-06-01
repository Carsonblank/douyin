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
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	//从token返回的登录用户Id
	var user_id int64
	//查看token是否合法
	if id, isValid := common.TokenValidity(token); !isValid {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Token is invalid.",
		},
		)
		//对于关注操作，token非法直接返回
		return
	} else {
		//合法的token解析出登陆用户ID并赋值给变量
		user_id = id
	}

	if _, exist := dtb.UserQueryByID(common.DB, user_id); exist {
		//获取关注对象的id
		to_user_id, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
		//查看关注对象id是否合法(是否存在)
		if _, exist := dtb.UserQueryByID(common.DB, to_user_id); !exist {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "To user doesn't exist"})
			return
		}
		//获取关注操作/取消关注操作的操作数
		action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
		//判断操作数是否合法
		if action_type == 1 {
			//判断是否已经关注了，防止重复关注
			if dtb.FavoriteQueryByUserAndUser(common.DB, user_id, to_user_id) {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Repeat follow."})
			} else {
				//数据库插入关注
				if err := dtb.RelationCreate(common.DB, dtb.Relation{
					UserId:   user_id,
					ToUserId: to_user_id,
				}); err != nil {
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: fmt.Sprintf("Relation create error: %v", err)})
				} else {
					//更新用户关注数和被关注数
					err0 := common.DB.Model(&dtb.User{}).Where("id = ?", user_id).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error
					err1 := common.DB.Model(&dtb.User{}).Where("id = ?", to_user_id).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error
					if err0 == nil && err1 == nil {
						c.JSON(http.StatusOK, Response{StatusCode: 0})
					} else {
						c.JSON(http.StatusOK, Response{
							StatusCode: 1, StatusMsg: fmt.Sprintf("Favorite action error : %v\n", err),
						})
					}
					c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Relation success"})
				}
			}

		} else if action_type == 2 {
			//判断本身有无关注，防止重复取消关注
			if !dtb.FavoriteQueryByUserAndUser(common.DB, user_id, to_user_id) {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Repeat delete follow."})
				return
			}

			if err := dtb.RelationDelete(common.DB, user_id, to_user_id); err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: fmt.Sprintf("Relation delete error: %v", err)})
			} else {
				//更新用户关注数和被关注数
				err0 := common.DB.Model(&dtb.User{}).Where("id = ?", user_id).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error
				err1 := common.DB.Model(&dtb.User{}).Where("id = ?", to_user_id).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error
				if err0 == nil && err1 == nil {
					c.JSON(http.StatusOK, Response{StatusCode: 0})
				} else {
					c.JSON(http.StatusOK, Response{
						StatusCode: 1, StatusMsg: fmt.Sprintf("Favorite action error : %v\n", err),
					})
				}
				c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Delete relation success\n"})
			}

		} else { //操作数不是1/2，操作非法
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Message action_type must be 1 or 2"})
		}

	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	//
	token := c.Query("token")
	var token_user_id int64
	//查看token是否合法
	if id, isValid := common.TokenValidity(token); !isValid {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Token is invalid.",
		},
		)
	} else {
		//合法的token解析出登陆用户ID并赋值给token_user_id变量
		token_user_id = id
		//查看合法token返回的用户ID是否存在
		if _, exist := dtb.UserQueryByID(common.DB, token_user_id); !exist {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
			)
		}
	}

	//
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if _, exist := dtb.UserQueryByID(common.DB, user_id); !exist {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User is not exist!\n",
			},
		})
		return
	}
	userList_db := dtb.UserQueryByFollowID(common.DB, user_id)

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: FromDBUsersTOMesUsers(common.DB, userList_db, token_user_id),
	})

}

// FollowerList all users have same follow list
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	//从token返回的登录用户Id
	var token_user_id int64
	//查看token是否合法
	if id, isValid := common.TokenValidity(token); !isValid {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Token is invalid.",
		},
		)
	} else {
		//合法的token解析出登陆用户ID并赋值给token_user_id变量
		token_user_id = id
		//查看合法token返回的用户ID是否存在
		if _, exist := dtb.UserQueryByID(common.DB, token_user_id); !exist {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
			)
		}
	}

	//获取需要查找粉丝列表的用户ID
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	//判断用户是否存在
	if _, exist := dtb.UserQueryByID(common.DB, user_id); !exist {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User is not exist!\n",
			},
		})
		return
	}
	userList_db := dtb.UserQueryByFollowerID(common.DB, user_id)

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: FromDBUsersTOMesUsers(common.DB, userList_db, token_user_id),
	})

}
