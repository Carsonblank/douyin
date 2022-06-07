package controller

import (
	"douyin/src/database"
	"douyin/src/repository"
	"douyin/src/service"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CommentAction(c *gin.Context) {
	//验证评论操作操作数是否合法
	actionType, err := service.CommentActionNum(c.Query("action_type"))
	if err != nil {
		//操作数不合法，返回失败响应
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("Comment action_type must be 1(add) or 2(del), you privode %d", actionType),
		})
		return
	}
	//验证token是否合法
	userId := service.Token2Id(c.Query("token"))
	if userId == 0 {
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  "Token is invalid.",
		},
		)
		return
	}
	//查看用户Id是否合法
	var userDb []database.User
	var exist bool
	if userDb, exist = repository.UserQueryByID(userId); !exist {
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		},
		)
		return
	}
	//查看videoId是否合法
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	//返回视频所有者Id
	var videoUserId int64
	if exist, id := repository.VideoQueryByID(videoId); !exist {
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1, StatusMsg: "Video is not exist!\n",
		})
		return
	} else {
		videoUserId = id
	}

	//评论操作actionType==1
	if actionType == 1 {
		//接收评论内容
		comment_text := c.Query("comment_text")
		//评论内容不能为空
		if comment_text == "" {
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  "Comment text must not null!\n",
			})
			return
		}
		//创建评论对象
		com := database.Comment{
			UserId:   userId,
			VideoId:  videoId,
			ToUserId: videoUserId,
			Content:  comment_text,
		}
		//添加到数据库
		if err := repository.CommentCreate(&com); err != nil {
			//添加失败
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("Comment action error : %v\n", err),
			})
			return
		}

		//转化为可用于报文的json格式,用于构建comment
		user := service.FromDBUsersTOMesUsers(userDb, userId)
		//构建comment对象,用于报文输出
		comment := service.Comment{
			Id:         com.Id,
			User:       user[0],
			Content:    com.Content,
			CreateDate: fmt.Sprintf("%d-%d", time.UnixMilli(com.CreatedAt).Month(), time.UnixMilli(com.CreatedAt).Day()),
		}
		//更新评论数
		if err := repository.CommentUpdataNumbers(videoId, true); err != nil {
			//更新失败
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1, StatusMsg: fmt.Sprintf("Comment action error : %v\n", err),
			})
			return
		}
		//返回成功响应
		c.JSON(http.StatusOK, service.CommentActionResponse{
			Response: service.Response{StatusCode: 0},
			Comment:  comment,
		})
		return
	}

	//删除评论操作actionType == 2
	if actionType == 2 {
		//获取评论Id
		commentId, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		//验证评论Id是否合法
		if !repository.CommentQueryByCommentId(commentId) {
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  "Comment is not exist.",
			})
		}
		//删除评论
		if err := repository.CommentDelete(commentId); err != nil {
			//删除失败
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("Delete comment action error : %v\n", err),
			})
		}
		//更新评论数
		if err := repository.CommentUpdataNumbers(videoId, false); err != nil {
			//更新失败
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1, StatusMsg: fmt.Sprintf("Comment action error : %v\n", err),
			})
			return
		}
		//返回成功响应
		c.JSON(http.StatusOK, service.CommentActionResponse{
			Response: service.Response{StatusCode: 0},
		})

	}
}
