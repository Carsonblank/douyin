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

//点赞操作
func FavoriteAction(db *gorm.DB) gin.HandlerFunc {
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
			video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
			action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
			//查看该视频是否存在
			if !dtb.VideoQueryByID(db, video_id) {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1, StatusMsg: "Video is not exist!\n",
				})
				return
			}
			//添加点赞记录(action_type==1)和去除点赞记录(action_type==2)
			if action_type == 1 {
				if exist := dtb.FavoriteQueryByUserAndVideo(db, user[0].Id, video_id); exist {
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Repeat favorite action."})
				} else {
					if err := dtb.FavoriteCreate(db, dtb.Favorite{
						UserId:   user[0].Id,
						ToUserId: dtb.UserQueryByVideoID(db, video_id),
						VideoId:  video_id,
					}); err != nil {
						c.JSON(http.StatusOK, Response{
							StatusCode: 1, StatusMsg: fmt.Sprintf("Favorite action error : %v\n", err),
						})
					} else {
						c.JSON(http.StatusOK, Response{StatusCode: 0})
					}
				}
			} else if action_type == 2 {
				//查看要取消的赞是否存在
				if exist := dtb.FavoriteQueryByUserAndVideo(db, user[0].Id, video_id); !exist {
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Repeat delete favorite action."})
					return
				}
				err := dtb.FavoriteDelete(db, user[0].Id, video_id)
				if err != nil { //取消赞失败
					c.JSON(http.StatusOK, Response{
						StatusCode: 1, StatusMsg: fmt.Sprintf("Remove favorite action error : %v\n", err),
					})
				} else { //取消赞成功
					c.JSON(http.StatusOK, Response{StatusCode: 0})
				}
			} else {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Favorite action_type must be 1 or 2"})
			}
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		}
	}
	return gin.HandlerFunc(fun)
}

//获取点赞列表
func FavoriteList(db *gorm.DB) gin.HandlerFunc {
	fun := func(c *gin.Context) {
		user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
		if user, exist := dtb.UserQueryByID(db, user_id); exist {
			videoList, _ := dtb.FavoriteQuerybyUserID(db, user_id)
			c.JSON(http.StatusOK, VideoListResponse{
				Response: Response{
					StatusCode: 0,
				},
				VideoList: FromDBVideosToMesVideos(db, videoList, user[0].Id),
			})
		} else {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: Response{
					StatusCode: 1,
				},
			})
		}

	}
	return gin.HandlerFunc(fun)
}
