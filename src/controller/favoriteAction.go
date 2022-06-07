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

//点赞操作
func FavoriteAction(c *gin.Context) {
	//验证点赞操作操作数是否合法
	actionType, err := service.FavoriteActionNum(c.Query("action_type"))
	if err != nil {
		//操作数不合法，返回失败响应
		c.JSON(http.StatusOK, service.Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("Favorite action_type must be 1(add) or 2(del), you privode %d", actionType),
		})
		return
	}
	//获取用户ID
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
	if _, exist := repository.UserQueryByID(userId); !exist {
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
	//赞操作
	if actionType == 1 {
		//验证添加赞操作是否合法，防止重复添加
		if exist := repository.FavoriteQueryByUserAndVideo(userId, videoId); exist {
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  "Repeat favorite action.",
			})
			return
		}

		//添加到数据库
		if err := repository.FavoriteCreate(database.Favorite{
			UserId:   userId,
			ToUserId: videoUserId,
			VideoId:  videoId,
		}); err != nil {
			//添加失败，返回失败响应
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("Favorite action error : %v\n", err),
			})
			return
		}
		//添加成功，统一更新视频点赞数、用户获赞数、用户喜欢数
		if err := repository.FavoriteUpdataNumbers(videoId, videoUserId, userId, true); err != nil {
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("Favorite action error : %v\n", err),
			})
			return
		}
		c.JSON(http.StatusOK, service.Response{StatusCode: 0})
		return
	}
	//取消赞操作
	if actionType == 2 {
		//验证取消赞操作是否合法，防止重复删除
		if exist := repository.FavoriteQueryByUserAndVideo(userId, videoId); !exist {
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  "Repeat favorite action.",
			})
			return
		}
		if err := repository.FavoriteDelete(userId, videoId); err != nil {
			//取消赞失败，返回失败响应
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1, StatusMsg: fmt.Sprintf("Remove favorite action error : %v\n", err),
			})
			return
		}
		//删除成功，统一更新视频点赞数、用户获赞数、用户喜欢数
		if err := repository.FavoriteUpdataNumbers(videoId, videoUserId, userId, false); err != nil {
			c.JSON(http.StatusOK, service.Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("Favorite action error : %v\n", err),
			})
			return
		}
		c.JSON(http.StatusOK, service.Response{StatusCode: 0})
		return
	}
}
