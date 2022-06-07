package controller

import (
	"douyin/src/database"
	"douyin/src/repository"
	"douyin/src/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RelationAction(c *gin.Context) {
	//验证action_type合法性
	actionType, err := service.RelationActionNum(c.Query("action_type"))
	if err != nil {
		//操作数不合法，返回失败响应
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("Relation action_type must be 1(add) or 2(del), you privode %d", actionType),
		})
		return
	}
	//验证token合法性
	userId := service.Token2Id(c.Query("token"))
	if userId == 0 {
		//返回失败响应
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  "Token is invalid.",
		},
		)
		return
	}
	//验证userId合法性
	if _, exist := repository.UserQueryByID(userId); !exist {
		//返回失败响应
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  "userId is invalid.",
		},
		)
		return
	}
	//验证to_user_id合法性
	toUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if _, exist := repository.UserQueryByID(toUserId); !exist {
		//返回失败响应
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  "to_user_id is invalid.",
		},
		)
		return
	}
	//关注操作
	if actionType == 1 {
		//验证关注操作是否合法（是否重复关注）
		if repository.RelationQueryByUserAndUser(userId, toUserId) {
			c.JSON(http.StatusOK, service.Response{StatusCode: 1, StatusMsg: "Repeat follow."})
			return
		}
		//数据库插入关注列表
		if err := repository.RelationCreate(database.Relation{
			UserId:   userId,
			ToUserId: toUserId,
		}); err != nil {
			//创建失败
			c.JSON(http.StatusOK,
				service.Response{
					StatusCode: 1,
					StatusMsg:  fmt.Sprintf("Relation create error: %v", err)})
			return
		}
		//更新用户关注数和被关注数
		if err := repository.RelationUpdataNumbers(userId, toUserId, true); err != nil {
			//更新失败，返回失败响应
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  err.Error()})
			return
		}
		//返回成功响应
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 0,
			StatusMsg:  "Relation success"})
		return

	}
	//取消关注操作
	if actionType == 2 {
		//验证取消关注操作是否合法（是否重复关注）
		if !repository.RelationQueryByUserAndUser(userId, toUserId) {
			c.JSON(http.StatusOK, service.Response{StatusCode: 1, StatusMsg: "Repeat delete follow."})
			return
		}
		//删除操作
		if err := repository.RelationDelete(userId, toUserId); err != nil {
			//返回删除失败响应
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("Relation delete error: %v", err)})
			return
		}
		//更新用户关注数和被关注数
		if err := repository.RelationUpdataNumbers(userId, toUserId, false); err != nil {
			//更新失败，返回失败响应
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  err.Error()})
			return
		}
		//返回成功响应
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 0,
			StatusMsg:  "Relation delete success"})
		return
	}
}
