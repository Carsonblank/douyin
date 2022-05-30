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
func CommentAction(db *gorm.DB) gin.HandlerFunc {
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

		if user_db, exist := dtb.UserQueryByToken(db, token); exist {
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
					ToUserId: dtb.UserQueryByVideoID(db, video_id),
					Content:  comment_text,
				}
				if err := dtb.CommentCreate(db, &com); err != nil {
					c.JSON(http.StatusOK, Response{
						StatusCode: 1,
						StatusMsg:  fmt.Sprintf("Comment action error : %v\n", err),
					})
				} else {
					user := FromDBUsersTOMesUsers(db, user_db, user_db[0].Id)
					comment := Comment{
						Id:         com.Id,
						User:       user[0],
						Content:    com.Content,
						CreateDate: fmt.Sprintf("%d-%d", time.UnixMilli(com.CreatedAt).Month(), time.UnixMilli(com.CreatedAt).Day()),
					}
					c.JSON(http.StatusOK, CommentActionResponse{
						Response: Response{StatusCode: 0},
						Comment:  comment,
					})
				}

			} else if action_type == 2 {
				comment_id, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)
				if err := dtb.CommentDelete(db, comment_id); err != nil {
					c.JSON(http.StatusOK, Response{
						StatusCode: 1,
						StatusMsg:  fmt.Sprintf("Delete comment action error : %v\n", err),
					})
				} else {
					c.JSON(http.StatusOK, Response{StatusCode: 0})
				}
			} else {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Comment action_type must be 1 or 2"})

			}
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		}
	}
	return gin.HandlerFunc(fun)
}

// CommentList all videos have same demo comment list
func CommentList(db *gorm.DB) gin.HandlerFunc {
	fun := func(c *gin.Context) {
		video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
		if !dtb.VideoQueryByID(db, video_id) {
			c.JSON(http.StatusOK, CommentListResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Video is not exist!\n"},
			})
			return
		}
		commentList, _ := dtb.CommentQuerybyVideoID(db, video_id)
		//获取登录用户的ID，用于是否关注判断
		token := c.Query("token")
		var userId int64
		if user, isfind := dtb.UserQueryByToken(db, token); isfind {
			userId = user[0].Id
		}
		//构建Messgae
		comments := make([]Comment, 0, len(commentList))
		for _, c := range commentList {
			user_db, _ := dtb.UserQueryByID(db, c.UserId)
			user := User{
				Id:            user_db[0].Id,
				Name:          user_db[0].Name,
				FollowCount:   dtb.UserFollowCount(db, user_db[0].Id),   //user_db.FollowCount,
				FollowerCount: dtb.UserFollowerCount(db, user_db[0].Id), //user_db.FollowerCount,
				IsFollow:      dtb.FavoriteQueryByUserAndUser(db, userId, user_db[0].Id),
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
	return gin.HandlerFunc(fun)
}
