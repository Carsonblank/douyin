//评论操作
package controller

import (
	"douyin/common"
	dtb "douyin/database"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

	token := c.Query("token")
	var token_user_id int64
	//查看token是否合法
	if id, isValid := common.TokenValidity(token); !isValid {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Token is invalid.",
		},
		)
		//对于评论操作，token非法直接返回
		return
	} else {
		//合法的token解析出登陆用户ID并赋值给token_user_id变量
		token_user_id = id
	}

	if user_db, exist := dtb.UserQueryByID(common.DB, token_user_id); exist {
		video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
		action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
		//发布评论(action_type==1)和删除评论(action_type==2)
		if action_type == 1 {
			comment_text := c.Query("comment_text")
			if comment_text == "" {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1,
					StatusMsg:  "Comment text must not null!\n",
				})
				return
			}
			com := dtb.Comment{
				UserId:   user_db[0].Id,
				VideoId:  video_id,
				ToUserId: dtb.UserQueryByVideoID(common.DB, video_id),
				Content:  comment_text,
			}
			if err := dtb.CommentCreate(common.DB, &com); err != nil {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1,
					StatusMsg:  fmt.Sprintf("Comment action error : %v\n", err),
				})
			} else {
				user := FromDBUsersTOMesUsers(common.DB, user_db, user_db[0].Id)
				comment := Comment{
					Id:         com.Id,
					User:       user[0],
					Content:    com.Content,
					CreateDate: fmt.Sprintf("%d-%d", time.UnixMilli(com.CreatedAt).Month(), time.UnixMilli(com.CreatedAt).Day()),
				}
				//更新评论数
				if err := common.DB.Model(&dtb.Video{}).Where("id = ?", video_id).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err == nil {
					c.JSON(http.StatusOK, CommentActionResponse{
						Response: Response{StatusCode: 0},
						Comment:  comment,
					})
				} else {
					c.JSON(http.StatusOK, Response{
						StatusCode: 1, StatusMsg: fmt.Sprintf("Comment action error : %v\n", err),
					})
				}
			}

		} else if action_type == 2 {
			comment_id, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)
			if err := dtb.CommentDelete(common.DB, comment_id); err != nil {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1,
					StatusMsg:  fmt.Sprintf("Delete comment action error : %v\n", err),
				})
			} else {
				//更新评论数
				if err := common.DB.Model(&dtb.Video{}).Where("id = ?", video_id).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err == nil {
					c.JSON(http.StatusOK, Response{StatusCode: 0})
				} else {
					c.JSON(http.StatusOK, Response{
						StatusCode: 1, StatusMsg: fmt.Sprintf("Comment action error : %v\n", err),
					})
				}
			}
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Comment action_type must be 1 or 2"})

		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if exist, _ := dtb.VideoQueryByID(common.DB, video_id); !exist {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Video is not exist!\n"},
		})
		return
	}
	commentList, _ := dtb.CommentQuerybyVideoID(common.DB, video_id)
	//获取登录用户的ID，用于是否关注判断
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
	//构建Messgae
	comments := make([]Comment, 0, len(commentList))
	for _, c := range commentList {
		user_db, _ := dtb.UserQueryByID(common.DB, c.UserId)
		user := User{
			Id:            user_db[0].Id,
			Name:          user_db[0].Name,
			FollowCount:   user_db[0].Follow_count,   //dtb.UserFollowCount(common.DB, user_db[0].Id),   //user_db.FollowCount,
			FollowerCount: user_db[0].Follower_count, //dtb.UserFollowerCount(common.DB, user_db[0].Id), //user_db.FollowerCount,
			IsFollow:      dtb.FavoriteQueryByUserAndUser(common.DB, token_user_id, user_db[0].Id),
			Avatar:        user_db[0].Avatar,
			Signature:     user_db[0].Signature,
		}
		comments = append(comments, Comment{
			Id:         c.Id,
			User:       user,
			Content:    c.Content,
			CreateDate: fmt.Sprintf("%d-%d", time.UnixMilli(c.CreatedAt).Month(), time.UnixMilli(c.CreatedAt).Day()),
		})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})

}
