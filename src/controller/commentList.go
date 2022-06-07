package controller

import (
	"douyin/src/repository"
	"douyin/src/service"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	//验证videoId是否合法
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if exist, _ := repository.VideoQueryByID(videoId); !exist {
		c.JSON(http.StatusOK, service.CommentListResponse{
			Response: service.Response{
				StatusCode: 1,
				StatusMsg:  "Video is not exist!\n"},
		})
		return
	}
	//获取评论列表
	commentList, _ := repository.CommentQuerybyVideoID(videoId)

	//验证token是否合法
	userId := service.Token2Id(c.Query("token"))

	//构建用于报文的Messgae
	comments := make([]service.Comment, 0, len(commentList))
	for _, c := range commentList {
		user_db, _ := repository.UserQueryByID(c.UserId)
		user := service.User{
			Id:            user_db[0].Id,
			Name:          user_db[0].Name,
			FollowCount:   user_db[0].Follow_count,
			FollowerCount: user_db[0].Follower_count,
			IsFollow:      repository.RelationQueryByUserAndUser(userId, user_db[0].Id),
			Avatar:        user_db[0].Avatar,
			Signature:     user_db[0].Signature,
		}
		comments = append(comments, service.Comment{
			Id:         c.Id,
			User:       user,
			Content:    c.Content,
			CreateDate: fmt.Sprintf("%d-%d", time.UnixMilli(c.CreatedAt).Month(), time.UnixMilli(c.CreatedAt).Day()),
		})
	}
	//发送成功响应
	c.JSON(http.StatusOK, service.CommentListResponse{
		Response:    service.Response{StatusCode: 0},
		CommentList: comments,
	})

}
